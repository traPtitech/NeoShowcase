package usecase

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/mattn/go-shellwords"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/samber/lo"

	"github.com/friendsofgo/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/util/progress/progressui"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	gstatus "google.golang.org/grpc/status"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

const (
	buildScriptName = "neoshowcase_internal_build.sh"
)

type BuilderService interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type builderService struct {
	client    domain.ControllerBuilderServiceClient
	buildkit  *buildkit.Client
	buildpack builder.BuildpackBackend
	storage   domain.Storage
	pubKey    *ssh.PublicKeys
	config    builder.ImageConfig

	appRepo      domain.ApplicationRepository
	artifactRepo domain.ArtifactRepository
	buildRepo    domain.BuildRepository
	gitRepo      domain.GitRepositoryRepository

	state       *state
	stateCancel func()
	statusLock  sync.Mutex
	response    chan<- *pb.BuilderResponse
	cancel      func()
}

func NewBuilderService(
	client domain.ControllerBuilderServiceClient,
	buildkit *buildkit.Client,
	buildpack builder.BuildpackBackend,
	storage domain.Storage,
	pubKey *ssh.PublicKeys,
	config builder.ImageConfig,
	appRepo domain.ApplicationRepository,
	artifactRepo domain.ArtifactRepository,
	buildRepo domain.BuildRepository,
	gitRepo domain.GitRepositoryRepository,
) BuilderService {
	return &builderService{
		client:       client,
		buildkit:     buildkit,
		buildpack:    buildpack,
		storage:      storage,
		pubKey:       pubKey,
		config:       config,
		appRepo:      appRepo,
		artifactRepo: artifactRepo,
		buildRepo:    buildRepo,
		gitRepo:      gitRepo,
	}
}

func (s *builderService) destImage(app *domain.Application, build *domain.Build) string {
	return s.config.FullImageName(app.ID) + ":" + build.Commit
}

func (s *builderService) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	response := make(chan *pb.BuilderResponse, 100)
	s.response = response

	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectBuilder(ctx, s.onRequest, response)
	}, 1*time.Second, 60*time.Second)
	go s.pruneLoop(ctx)

	return nil
}

func (s *builderService) Shutdown(_ context.Context) error {
	s.cancel()
	s.statusLock.Lock()
	defer s.statusLock.Unlock()
	if s.stateCancel != nil {
		s.stateCancel()
	}
	return nil
}

func (s *builderService) pruneLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.prune(ctx)
			if err != nil {
				log.Errorf("failed to prune buildkit: %+v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *builderService) prune(ctx context.Context) error {
	return s.buildkit.Prune(ctx, nil, buildkit.PruneAll)
}

func (s *builderService) cancelBuild(buildID string) {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil && s.stateCancel != nil {
		if s.state.build.ID == buildID {
			s.stateCancel()
		}
	}
}

func (s *builderService) authSessions() []session.Attachable {
	if s.config.Registry.Username == "" && s.config.Registry.Password == "" {
		return nil
	}
	return []session.Attachable{authprovider.NewDockerAuthProvider(&configfile.ConfigFile{
		AuthConfigs: map[string]types.AuthConfig{
			s.config.Registry.Addr: {
				Username: s.config.Registry.Username,
				Password: s.config.Registry.Password,
			},
		},
	})}
}

func (s *builderService) onRequest(req *pb.BuilderRequest) {
	switch req.Type {
	case pb.BuilderRequest_START_BUILD:
		b := req.Body.(*pb.BuilderRequest_StartBuild).StartBuild
		err := s.tryStartTask(b.BuildId)
		if err != nil {
			log.Errorf("failed to start build: %+v", err)
		}
	case pb.BuilderRequest_CANCEL_BUILD:
		b := req.Body.(*pb.BuilderRequest_CancelBuild).CancelBuild
		s.cancelBuild(b.BuildId)
	default:
		log.Errorf("unknown builder request type: %v", req.Type)
	}
}

func (s *builderService) tryStartTask(buildID string) error {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil {
		log.Infof("skipping build request for %v, builder busy", buildID)
		return nil // Builder busy - skip
	}

	now := time.Now()
	err := s.buildRepo.UpdateBuild(context.Background(), buildID, domain.UpdateBuildArgs{
		FromStatus: optional.From(domain.BuildStatusQueued),
		Status:     optional.From(domain.BuildStatusBuilding),
		StartedAt:  optional.From(now),
		UpdatedAt:  optional.From(now),
	})
	if err == repository.ErrNotFound {
		return nil // other builder has acquired the build lock - skip
	}
	if err != nil {
		return err
	}

	// Acquired build lock
	log.Infof("Starting build for %v", buildID)

	build, err := s.buildRepo.GetBuild(context.Background(), buildID)
	if err != nil {
		return err
	}
	app, err := s.appRepo.GetApplication(context.Background(), build.ApplicationID)
	if err != nil {
		return err
	}
	repo, err := s.gitRepo.GetRepository(context.Background(), app.RepositoryID)
	if err != nil {
		return err
	}

	st, err := newState(app, build, repo, s.response)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	finishWait := make(chan struct{})
	s.state = st
	s.stateCancel = func() {
		cancel()
		<-finishWait
	}

	go func() {
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_STARTED, Body: &pb.BuilderResponse_Started{Started: &pb.BuildStarted{
			BuildId: buildID,
		}}}

		status := s.process(ctx, st)
		s.finalize(context.Background(), st, status) // don't want finalization tasks to be cancelled
		st.Cleanup()

		cancel()
		close(finishWait)
		s.statusLock.Lock()
		s.state = nil
		s.stateCancel = nil
		s.statusLock.Unlock()
		log.Infof("Build settled for %v", buildID)
		// Send settled response *after* unlocking internal state for next build
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_SETTLED, Body: &pb.BuilderResponse_Settled{Settled: &pb.BuildSettled{
			BuildId: buildID,
			Reason:  toPBSettleReason(status),
		}}}
	}()

	return nil
}

func (s *builderService) finalize(ctx context.Context, st *state, status domain.BuildStatus) {
	err := domain.SaveBuildLog(s.storage, st.build.ID, st.logWriter.LogReader())
	if err != nil {
		log.Errorf("failed to save build log: %+v", err)
	}

	now := time.Now()
	updateArgs := domain.UpdateBuildArgs{
		Status:     optional.From(status),
		UpdatedAt:  optional.From(now),
		FinishedAt: optional.From(now),
	}
	if err := s.buildRepo.UpdateBuild(ctx, st.build.ID, updateArgs); err != nil {
		log.Errorf("failed to update build: %+v", err)
	}
}

type buildStep struct {
	desc string
	fn   func() error
}

func (s *builderService) buildSteps(ctx context.Context, st *state) ([]buildStep, error) {
	var steps []buildStep

	steps = append(steps, buildStep{"Repository Clone", func() error {
		return s.cloneRepository(ctx, st)
	}})

	switch bc := st.app.Config.BuildConfig.(type) {
	case *domain.BuildConfigRuntimeBuildpack:
		steps = append(steps, buildStep{"Build (Runtime Buildpack)", func() error {
			return s.buildImageBuildpack(ctx, st, bc)
		}})
	case *domain.BuildConfigRuntimeCmd:
		steps = append(steps, buildStep{"Build (Runtime Command)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildImageWithCmd(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigRuntimeDockerfile:
		steps = append(steps, buildStep{"Build (Runtime Dockerfile)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildImageWithDockerfile(ctx, st, ch, bc)
			})
		}})
	case *domain.BuildConfigStaticCmd:
		steps = append(steps, buildStep{"Build (Static Command)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticWithCmd(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Save Artifact", func() error {
			return s.saveArtifact(ctx, st)
		}})
	case *domain.BuildConfigStaticDockerfile:
		steps = append(steps, buildStep{"Build (Static Dockerfile)", func() error {
			return withBuildkitProgress(ctx, st.logWriter, func(ctx context.Context, ch chan *buildkit.SolveStatus) error {
				return s.buildStaticWithDockerfile(ctx, st, ch, bc)
			})
		}})
		steps = append(steps, buildStep{"Save Artifact", func() error {
			return s.saveArtifact(ctx, st)
		}})
	default:
		return nil, errors.New("unknown build config type")
	}

	return steps, nil
}

func (s *builderService) process(ctx context.Context, st *state) domain.BuildStatus {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go s.updateStatusLoop(ctx, st.build.ID)

	steps, err := s.buildSteps(ctx, st)
	if err != nil {
		log.Errorf("calculating build steps: %+v", err)
		st.WriteLog(fmt.Sprintf("[ns-builder] Error calculating build steps: %v", err))
		return domain.BuildStatusFailed
	}

	for i, step := range steps {
		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) %s", i+1, len(steps), step.desc))
		start := time.Now()
		err := step.fn()
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, gstatus.FromContextError(context.Canceled).Err()) {
			st.WriteLog("[ns-builder] Build cancelled.")
			return domain.BuildStatusCanceled
		}
		if err != nil {
			msg := fmt.Sprintf("%+v", err)
			log.Error(msg)
			st.WriteLog("[ns-builder] Build failed:")
			st.WriteLog(msg)
			return domain.BuildStatusFailed
		}
		st.WriteLog(fmt.Sprintf("[ns-builder] ==> (%d/%d) Done (%v).", i+1, len(steps), time.Since(start)))
	}

	return domain.BuildStatusSucceeded
}

func (s *builderService) updateStatusLoop(ctx context.Context, buildID string) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := s.buildRepo.UpdateBuild(ctx, buildID, domain.UpdateBuildArgs{UpdatedAt: optional.From(time.Now())})
			if err != nil {
				log.Errorf("failed to update build time: %+v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *builderService) cloneRepository(ctx context.Context, st *state) error {
	repo, err := git.PlainInit(st.repositoryTempDir, false)
	if err != nil {
		return errors.Wrap(err, "failed to init repository")
	}
	auth, err := domain.GitAuthMethod(st.repo, s.pubKey)
	if err != nil {
		return err
	}
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{st.repo.URL},
	})
	if err != nil {
		return errors.Wrap(err, "failed to add remote")
	}
	targetRef := plumbing.NewRemoteReferenceName("origin", "target")
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", st.build.Commit, targetRef))},
		Depth:      1,
		Auth:       auth,
	})
	if err != nil {
		return errors.Wrap(err, "failed to clone repository")
	}
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}
	err = wt.Checkout(&git.CheckoutOptions{Branch: targetRef})
	if err != nil {
		return errors.Wrap(err, "failed to checkout")
	}
	sm, err := wt.Submodules()
	if err != nil {
		return errors.Wrap(err, "getting submodules")
	}
	err = sm.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})
	if err != nil {
		return errors.Wrap(err, "updating submodules")
	}
	return nil
}

func (s *builderService) saveArtifact(ctx context.Context, st *state) error {
	filename := st.artifactTempFile.Name()

	stat, err := os.Stat(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}

	artifact := domain.NewArtifact(st.build.ID, stat.Size())
	err = s.artifactRepo.CreateArtifact(ctx, artifact)
	if err != nil {
		return errors.Wrap(err, "creating artifact record")
	}

	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "opening artifact")
	}
	defer file.Close()

	pr, pw := io.Pipe()
	gzipWriter := gzip.NewWriter(pw)
	if err != nil {
		return errors.Wrap(err, "creating gzip stream")
	}
	go func() {
		defer pw.Close()
		_, err := io.Copy(gzipWriter, file)
		if err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "copying file to pipe writer"))
			return
		}
		err = gzipWriter.Close()
		if err != nil {
			_ = pw.CloseWithError(errors.Wrap(err, "flushing gzip writer"))
			return
		}
	}()
	err = domain.SaveArtifact(s.storage, artifact.ID, pr)
	if err != nil {
		return errors.Wrap(err, "saving artifact")
	}

	return nil
}

func withBuildkitProgress(ctx context.Context, logger io.Writer, buildFn func(ctx context.Context, ch chan *buildkit.SolveStatus) error) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return buildFn(ctx, ch)
	})
	eg.Go(func() error {
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, logger, ch)
		return err
	})
	return eg.Wait()
}

func (s *builderService) buildImageBuildpack(
	ctx context.Context,
	st *state,
	bc *domain.BuildConfigRuntimeBuildpack,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	buildDir := filepath.Join(st.repositoryTempDir, contextDir)
	return s.buildpack.Pack(ctx, buildDir, st.Logger(), s.destImage(st.app, st.build))
}

func (s *builderService) buildImageWithCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeCmd,
) error {
	var ls llb.State
	if bc.BaseImage == "" {
		ls = llb.Scratch()
	} else {
		ls = llb.Image(bc.BaseImage)
	}
	ls = ls.
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			AllowWildcard:       true,
			CreateDestPath:      true,
		}))

	if bc.BuildCmd != "" {
		if bc.BuildCmdShell {
			err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args([]string{"./" + buildScriptName})).Root()
		} else {
			args, err := shellwords.Parse(bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args(args)).Root()
		}
	}

	def, err := ls.Marshal(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to marshal llb")
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": s.destImage(st.app, st.build),
				"push": "true",
			},
		}},
		LocalDirs: map[string]string{
			"local-src": st.repositoryTempDir,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildImageWithDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigRuntimeDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, ".")
	dockerfileDir := filepath.Join(contextDir, filepath.Dir(bc.DockerfileName))
	_, err := s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type: buildkit.ExporterImage,
			Attrs: map[string]string{
				"name": s.destImage(st.app, st.build),
				"push": "true",
			},
		}},
		LocalDirs: map[string]string{
			"context":    filepath.Join(st.repositoryTempDir, contextDir),
			"dockerfile": filepath.Join(st.repositoryTempDir, dockerfileDir),
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": bc.DockerfileName},
		Session:       s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildStaticWithCmd(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticCmd,
) error {
	var ls llb.State
	if bc.BaseImage == "" {
		ls = llb.Scratch()
	} else {
		ls = llb.Image(bc.BaseImage)
	}
	ls = ls.
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			AllowWildcard:       true,
			CreateDestPath:      true,
		}))

	if bc.BuildCmd != "" {
		if bc.BuildCmdShell {
			err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args([]string{"./" + buildScriptName})).Root()
		} else {
			args, err := shellwords.Parse(bc.BuildCmd)
			if err != nil {
				return err
			}
			ls = ls.Run(llb.Args(args)).Root()
		}
	}

	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(ls, filepath.Join("/srv", bc.ArtifactPath), "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to marshal llb")
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"local-src": st.repositoryTempDir,
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func (s *builderService) buildStaticWithDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
	bc *domain.BuildConfigStaticDockerfile,
) error {
	contextDir := lo.Ternary(bc.Context != "", bc.Context, filepath.Dir(bc.DockerfileName))
	dockerfile, err := os.ReadFile(filepath.Join(st.repositoryTempDir, contextDir, bc.DockerfileName))
	if err != nil {
		return err
	}

	b, _, _, err := dockerfile2llb.Dockerfile2LLB(ctx, dockerfile, dockerfile2llb.ConvertOpt{})
	if err != nil {
		return err
	}

	def, err := llb.
		Scratch().
		File(llb.Copy(*b, bc.ArtifactPath, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return err
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"context": filepath.Join(st.repositoryTempDir, contextDir),
		},
		Session: s.authSessions(),
	}, ch)
	return err
}

func createScriptFile(filename string, script string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("#!/bin/sh\nset -eux\n" + script)
	if err != nil {
		return err
	}
	return nil
}

type state struct {
	app       *domain.Application
	build     *domain.Build
	repo      *domain.Repository
	logWriter *logWriter

	repositoryTempDir string
	artifactTempFile  *os.File
}

type logWriter struct {
	buildID  string
	response chan<- *pb.BuilderResponse
	buf      bytes.Buffer
}

func (l *logWriter) toBuilderResponse(p []byte) *pb.BuilderResponse {
	return &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_LOG, Body: &pb.BuilderResponse_Log{
		Log: &pb.BuildLogPortion{BuildId: l.buildID, Log: p},
	}}
}

func (l *logWriter) LogReader() io.Reader {
	return &l.buf
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	n, err = l.buf.Write(p)
	if err != nil {
		return
	}
	select {
	case l.response <- l.toBuilderResponse(p):
	default:
	}
	return
}

func newState(app *domain.Application, build *domain.Build, repo *domain.Repository, response chan<- *pb.BuilderResponse) (*state, error) {
	st := &state{
		app:   app,
		build: build,
		repo:  repo,
		logWriter: &logWriter{
			buildID:  build.ID,
			response: response,
		},
	}
	var err error
	st.repositoryTempDir, err = os.MkdirTemp("", "repository-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp repository dir")
	}
	st.artifactTempFile, err = os.CreateTemp("", "artifacts-")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp artifact file")
	}
	return st, nil
}

func (s *state) Cleanup() {
	_ = os.RemoveAll(s.repositoryTempDir)
	_ = os.Remove(s.artifactTempFile.Name())
}

func (s *state) Logger() io.Writer {
	return s.logWriter
}

func (s *state) WriteLog(a ...interface{}) {
	_, _ = fmt.Fprintln(s.logWriter, a...)
}

func toPBSettleReason(status domain.BuildStatus) pb.BuildSettled_Reason {
	switch status {
	case domain.BuildStatusSucceeded:
		return pb.BuildSettled_SUCCESS
	case domain.BuildStatusFailed:
		return pb.BuildSettled_FAILED
	case domain.BuildStatusCanceled:
		return pb.BuildSettled_CANCELLED
	default:
		panic(fmt.Sprintf("unexpected settled status: %v", status))
	}
}
