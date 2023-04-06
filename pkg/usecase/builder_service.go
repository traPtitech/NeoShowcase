package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
	"github.com/traPtitech/neoshowcase/pkg/util/retry"
)

const (
	buildScriptName      = "neoshowcase_internal_build.sh"
	entryPointScriptName = "neoshowcase_internal_entrypoint.sh"
)

type BuilderService interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type builderService struct {
	client   domain.ComponentServiceClient
	buildkit *buildkit.Client
	storage  domain.Storage
	pubKey   *ssh.PublicKeys

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
	client domain.ComponentServiceClient,
	buildkit *buildkit.Client,
	storage domain.Storage,
	pubKey *ssh.PublicKeys,
	artifactRepo domain.ArtifactRepository,
	buildRepo domain.BuildRepository,
	gitRepo domain.GitRepositoryRepository,
) BuilderService {
	return &builderService{
		client:       client,
		buildkit:     buildkit,
		storage:      storage,
		pubKey:       pubKey,
		artifactRepo: artifactRepo,
		buildRepo:    buildRepo,
		gitRepo:      gitRepo,
	}
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
		if s.state.task.BuildID == buildID {
			s.stateCancel()
		}
	}
}

func (s *builderService) onRequest(req *pb.BuilderRequest) {
	switch req.Type {
	case pb.BuilderRequest_START_BUILD_IMAGE:
		b := req.Body.(*pb.BuilderRequest_BuildImage).BuildImage
		err := s.tryStartTask(&builder.Task{
			BuildID:       b.BuildId,
			ApplicationID: b.ApplicationId,
			Static:        false,
			BuildSource: &builder.BuildSource{
				RepositoryID: b.Source.RepositoryId,
				Commit:       b.Source.Commit,
			},
			BuildOptions: &builder.BuildOptions{
				BaseImageName:  b.Options.BaseImageName,
				DockerfileName: b.Options.DockerfileName,
				ArtifactPath:   b.Options.ArtifactPath,
				BuildCmd:       b.Options.BuildCmd,
				EntrypointCmd:  b.Options.EntrypointCmd,
			},
			ImageName: b.ImageName,
			ImageTag:  b.ImageTag,
		})
		if err != nil {
			log.Errorf("failed to start build: %+v", err)
		}
	case pb.BuilderRequest_START_BUILD_STATIC:
		b := req.Body.(*pb.BuilderRequest_BuildStatic).BuildStatic
		err := s.tryStartTask(&builder.Task{
			BuildID:       b.BuildId,
			ApplicationID: b.ApplicationId,
			Static:        true,
			BuildSource: &builder.BuildSource{
				RepositoryID: b.Source.RepositoryId,
				Commit:       b.Source.Commit,
			},
			BuildOptions: &builder.BuildOptions{
				BaseImageName:  b.Options.BaseImageName,
				DockerfileName: b.Options.DockerfileName,
				ArtifactPath:   b.Options.ArtifactPath,
				BuildCmd:       b.Options.BuildCmd,
				EntrypointCmd:  b.Options.EntrypointCmd,
			},
		})
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

func (s *builderService) tryStartTask(task *builder.Task) error {
	s.statusLock.Lock()
	defer s.statusLock.Unlock()

	if s.state != nil {
		return errors.New("builder unavailable")
	}

	now := time.Now()
	err := s.buildRepo.UpdateBuild(context.Background(), task.BuildID, domain.UpdateBuildArgs{
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

	log.Infof("Starting build for %v", task.BuildID)

	repo, err := s.gitRepo.GetRepository(context.Background(), task.BuildSource.RepositoryID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	finish := make(chan struct{})
	st := newState(task, repo, s.response)
	s.state = st
	s.stateCancel = func() {
		cancel()
		<-finish
	}

	go func() {
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_STARTED, Body: &pb.BuilderResponse_Started{Started: &pb.BuildStarted{
			ApplicationId: task.ApplicationID,
			BuildId:       task.BuildID,
		}}}
		status := s.process(ctx, st)
		s.finalize(context.Background(), st, status) // don't want finalization tasks to be cancelled
		s.response <- &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_SETTLED, Body: &pb.BuilderResponse_Settled{Settled: &pb.BuildSettled{
			ApplicationId: task.ApplicationID,
			BuildId:       task.BuildID,
			Reason:        toPBSettleReason(status),
		}}}

		cancel()
		close(finish)
		s.statusLock.Lock()
		s.state = nil
		s.stateCancel = nil
		s.statusLock.Unlock()
		log.Infof("Build settled for %v", task.BuildID)
	}()

	return nil
}

func (s *builderService) process(ctx context.Context, st *state) domain.BuildStatus {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go s.updateStatusLoop(ctx, st.task.BuildID)

	err := st.initTempFiles(st.task.Static)
	if err != nil {
		log.Errorf("failed to init temp files: %+v", err)
		return domain.BuildStatusFailed
	}

	err = s.cloneRepository(ctx, st)
	if err != nil {
		log.Errorf("failed to clone repository: %+v", err)
		return domain.BuildStatusFailed
	}

	return s.build(ctx, st)
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
	auth, err := domain.GitAuthMethod(st.repository, s.pubKey)
	if err != nil {
		return err
	}
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{st.repository.URL},
	})
	if err != nil {
		return errors.Wrap(err, "failed to add remote")
	}
	targetRef := plumbing.NewRemoteReferenceName("origin", "target")
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", st.task.BuildSource.Commit, targetRef))},
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
	return nil
}

func (s *builderService) build(ctx context.Context, st *state) domain.BuildStatus {
	st.writeLog("[ns-builder] Build started.")

	var err error
	if st.task.Static {
		err = s.buildStatic(ctx, st)
	} else {
		err = s.buildImage(ctx, st)
	}
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded || errors.Is(err, gstatus.FromContextError(context.Canceled).Err()) {
			st.writeLog("[ns-builder] Build cancelled.")
			return domain.BuildStatusCanceled
		}
		log.Errorf("failed to build: %+v", err)
		return domain.BuildStatusFailed
	}

	st.writeLog("[ns-builder] Build succeeded!")
	return domain.BuildStatusSucceeded
}

func (s *builderService) finalize(ctx context.Context, st *state, status domain.BuildStatus) {
	// ログファイルの保存
	if st.logTempFile != nil {
		_ = st.logTempFile.Close()
		if err := domain.SaveBuildLog(s.storage, st.logTempFile.Name(), st.task.BuildID); err != nil {
			log.Errorf("failed to save build log: %+v", err)
		}
	}

	// 生成物tarの保存
	if st.artifactTempFile != nil {
		_ = st.artifactTempFile.Close()
		if status == domain.BuildStatusSucceeded {
			err := func() error {
				filename := st.artifactTempFile.Name()
				stat, err := os.Stat(filename)
				if err != nil {
					return errors.Wrap(err, "failed to open artifact")
				}
				artifact := domain.NewArtifact(st.task.BuildID, stat.Size())
				err = s.artifactRepo.CreateArtifact(ctx, artifact)
				if err != nil {
					return errors.Wrap(err, "failed to create artifact")
				}
				err = domain.SaveArtifact(s.storage, filename, artifact.ID)
				if err != nil {
					return errors.Wrap(err, "failed to save artifact")
				}
				return nil
			}()
			log.Errorf("failed to process artifact: %+v", err)
		} else {
			_ = os.Remove(st.artifactTempFile.Name())
		}
	}

	// 一時リポジトリディレクトリの削除
	if st.repositoryTempDir != "" {
		_ = os.RemoveAll(st.repositoryTempDir)
	}

	// Build更新
	now := time.Now()
	updateArgs := domain.UpdateBuildArgs{
		Status:     optional.From(status),
		UpdatedAt:  optional.From(now),
		FinishedAt: optional.From(now),
	}
	if err := s.buildRepo.UpdateBuild(ctx, st.task.BuildID, updateArgs); err != nil {
		log.Errorf("failed to update build: %+v", err)
	}
}

func (s *builderService) buildImage(ctx context.Context, st *state) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		// イメージの出力先設定
		exportAttrs := map[string]string{
			"name": st.task.ImageName + ":" + st.task.ImageTag,
			"push": "true",
		}
		if st.task.BuildOptions.DockerfileName != "" {
			// Dockerfileを使用
			// entrypoint, startupコマンドは無視
			err = s.buildImageWithDockerfile(ctx, st, exportAttrs, ch)
		} else {
			// 指定したベースイメージを使用
			err = s.buildImageWithConfig(ctx, st, exportAttrs, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, st.getLogWriter(), ch)
		return err
	})
	return eg.Wait()
}

func (s *builderService) buildImageWithDockerfile(
	ctx context.Context,
	st *state,
	exportAttrs map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	contextDir := filepath.Dir(st.task.BuildOptions.DockerfileName)
	_, err := s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:  buildkit.ExporterImage,
			Attrs: exportAttrs,
		}},
		LocalDirs: map[string]string{
			"context":    filepath.Join(st.repositoryTempDir, contextDir),
			"dockerfile": st.repositoryTempDir,
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": st.task.BuildOptions.DockerfileName},
	}, ch)
	return err
}

func (s *builderService) buildImageWithConfig(
	ctx context.Context,
	st *state,
	exportAttrs map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), st.task.BuildOptions.BuildCmd)
	if err != nil {
		return err
	}

	err = createScriptFile(filepath.Join(st.repositoryTempDir, entryPointScriptName), st.task.BuildOptions.EntrypointCmd)
	if err != nil {
		return err
	}

	dockerfile := fmt.Sprintf(
		`FROM %s
WORKDIR /srv
COPY . .
RUN ["/srv/%s"]
ENTRYPOINT ["/srv/%s"]`,
		util.ValueOr(st.task.BuildOptions.BaseImageName, "scratch"),
		buildScriptName,
		entryPointScriptName,
	)
	b, _, _, err := dockerfile2llb.Dockerfile2LLB(ctx, []byte(dockerfile), dockerfile2llb.ConvertOpt{})
	if err != nil {
		return err
	}

	def, err := b.Marshal(ctx)
	if err != nil {
		return err
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:  buildkit.ExporterImage,
			Attrs: exportAttrs,
		}},
		LocalDirs: map[string]string{
			"context": st.repositoryTempDir,
		},
		Frontend: "dockerfile.v0",
	}, ch)
	return err
}

func (s *builderService) buildStatic(ctx context.Context, st *state) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if st.task.BuildOptions.DockerfileName != "" {
			// Dockerfileを使用
			// entrypoint, startupコマンドは無視
			err = s.buildStaticWithDockerfile(ctx, st, ch)
		} else {
			// 指定したベースイメージを使用
			err = s.buildStaticWithConfig(ctx, st, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, st.getLogWriter(), ch)
		return err
	})
	return eg.Wait()
}

func (s *builderService) buildStaticWithDockerfile(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	dockerfile, err := os.ReadFile(filepath.Join(st.repositoryTempDir, st.task.BuildOptions.DockerfileName))
	if err != nil {
		return err
	}

	b, _, _, err := dockerfile2llb.Dockerfile2LLB(ctx, dockerfile, dockerfile2llb.ConvertOpt{})
	if err != nil {
		return err
	}

	def, err := llb.
		Scratch().
		File(llb.Copy(*b, st.task.BuildOptions.ArtifactPath, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(ctx)
	if err != nil {
		return err
	}

	contextDir := filepath.Dir(st.task.BuildOptions.DockerfileName)
	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"context": filepath.Join(st.repositoryTempDir, contextDir),
		},
	}, ch)
	return err
}

func (s *builderService) buildStaticWithConfig(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	err := createScriptFile(filepath.Join(st.repositoryTempDir, buildScriptName), st.task.BuildOptions.BuildCmd)
	if err != nil {
		return err
	}

	var ls llb.State
	if st.task.BuildOptions.BaseImageName == "" {
		ls = llb.Scratch()
	} else {
		ls = llb.Image(st.task.BuildOptions.BaseImageName)
	}
	b := ls.
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			AllowWildcard:       true,
			CreateDestPath:      true,
		})).
		Run(llb.Args([]string{"./" + buildScriptName})).
		Root()

	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(b, filepath.Join("/srv", st.task.BuildOptions.ArtifactPath), "/", &llb.CopyInfo{
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
			"local-src": st.repositoryTempDir,
		},
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
	task       *builder.Task
	repository *domain.Repository
	response   chan<- *pb.BuilderResponse

	repositoryTempDir string
	logTempFile       *os.File
	logWriter         *logWriter
	artifactTempFile  *os.File
}

type logWriter struct {
	buildID  string
	response chan<- *pb.BuilderResponse
	logFile  *os.File
}

func (l *logWriter) toBuilderResponse(p []byte) *pb.BuilderResponse {
	return &pb.BuilderResponse{Type: pb.BuilderResponse_BUILD_LOG, Body: &pb.BuilderResponse_Log{
		Log: &pb.BuildLogPortion{BuildId: l.buildID, Log: p},
	}}
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	n, err = l.logFile.Write(p)
	if err != nil {
		return
	}
	select {
	case l.response <- l.toBuilderResponse(p):
	default:
	}
	return
}

func newState(task *builder.Task, repo *domain.Repository, response chan<- *pb.BuilderResponse) *state {
	return &state{
		task:       task,
		repository: repo,
		response:   response,
	}
}

func (s *state) initTempFiles(useArtifactTempFile bool) error {
	var err error

	// ログ用一時ファイル作成
	s.logTempFile, err = os.CreateTemp("", "buildlog-")
	if err != nil {
		return errors.Wrap(err, "failed to create tmp log file")
	}
	s.logWriter = &logWriter{
		buildID:  s.task.BuildID,
		response: s.response,
		logFile:  s.logTempFile,
	}

	// 成果物tarの一時保存先作成
	if useArtifactTempFile {
		s.artifactTempFile, err = os.CreateTemp("", "artifacts-")
		if err != nil {
			return errors.Wrap(err, "failed to create tmp artifact file")
		}
	}

	// リポジトリクローン用の一時ディレクトリ作成
	s.repositoryTempDir, err = os.MkdirTemp("", "repository-")
	if err != nil {
		return errors.Wrap(err, "failed to create tmp repository dir")
	}

	return nil
}

func (s *state) getLogWriter() io.Writer {
	return s.logWriter
}

func (s *state) writeLog(a ...interface{}) {
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
