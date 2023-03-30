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
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	gstatus "google.golang.org/grpc/status"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
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

	// TODO: 後で消す
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

	response := make(chan *pb.BuilderResponse)
	s.response = response

	go retry.Do(ctx, func(ctx context.Context) error {
		return s.client.ConnectBuilder(ctx, s.onRequest, response)
	}, 1*time.Second, 60*time.Second)

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
			log.WithError(err).Errorf("failed to start build: %v", err)
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
			log.WithError(err).Errorf("failed to start build: %v", err)
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
		FromStatus: optional.From(builder.BuildStatusQueued),
		Status:     optional.From(builder.BuildStatusBuilding),
		StartedAt:  optional.From(now),
		UpdatedAt:  optional.From(now),
	})
	if err == repository.ErrNotFound {
		return nil // other builder has acquired the build lock - skip
	}
	if err != nil {
		return err
	}

	repo, err := s.gitRepo.GetRepository(context.Background(), task.BuildID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	finish := make(chan struct{})
	st := newState(task, repo)
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
	}()

	return nil
}

func (s *builderService) process(ctx context.Context, st *state) builder.BuildStatus {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go s.updateStatusLoop(ctx, st.task.BuildID)

	err := st.initTempFiles(st.task.Static)
	if err != nil {
		log.WithError(err).Error("failed to init temp files")
		return builder.BuildStatusFailed
	}

	err = s.cloneRepository(ctx, st)
	if err != nil {
		log.WithError(err).Error("failed to clone repository")
		return builder.BuildStatusFailed
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
				log.WithError(err).Error("failed to update build time")
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

func (s *builderService) build(ctx context.Context, st *state) builder.BuildStatus {
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
			return builder.BuildStatusCanceled
		}
		log.WithError(err).Error("failed to build")
		return builder.BuildStatusFailed
	}

	st.writeLog("[ns-builder] Build succeeded!")
	return builder.BuildStatusSucceeded
}

func (s *builderService) finalize(ctx context.Context, st *state, status builder.BuildStatus) {
	// ログファイルの保存
	if st.logTempFile != nil {
		_ = st.logTempFile.Close()
		if err := domain.SaveLogFile(s.storage, st.logTempFile.Name(), filepath.Join("buildlogs", st.task.BuildID), st.task.BuildID); err != nil {
			log.WithError(err).Error("failed to save build log")
		}
	}

	// 生成物tarの保存
	if st.artifactTempFile != nil {
		_ = st.artifactTempFile.Close()
		if status == builder.BuildStatusSucceeded {
			sid := domain.NewID()
			filename := st.artifactTempFile.Name()
			stat, err := os.Stat(filename)
			if err != nil {
				log.WithError(err).Error("failed to open artifact")
			} else {
				err = s.artifactRepo.CreateArtifact(ctx, stat.Size(), st.task.BuildID, sid)
				if err != nil {
					log.WithError(err).Error("failed to create artifact")
				}
			}

			err = domain.SaveArtifact(s.storage, filename, filepath.Join("artifacts", fmt.Sprintf("%s.tar", sid)))
			if err != nil {
				log.WithError(err).Error("failed to save directory to tar")
			}
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
		log.WithError(err).Error("failed to update build_log entry")
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
		if len(st.task.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
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
	_, err := s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:  buildkit.ExporterImage,
			Attrs: exportAttrs,
		}},
		LocalDirs: map[string]string{
			"context":    st.repositoryTempDir,
			"dockerfile": st.repositoryTempDir,
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": st.task.BuildOptions.DockerfileName},
		Session:       []session.Attachable{authprovider.NewDockerAuthProvider(io.Discard)},
	}, ch)
	return err
}

func (s *builderService) buildImageWithConfig(
	ctx context.Context,
	st *state,
	exportAttrs map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	var fs, fe *os.File
	fs, err := os.OpenFile(filepath.Join(st.repositoryTempDir, buildScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fs.WriteString("#!/bin/sh\n" + st.task.BuildOptions.BuildCmd)
	if err != nil {
		return err
	}
	defer fs.Close()
	defer os.Remove(fs.Name())

	fe, err = os.OpenFile(filepath.Join(st.repositoryTempDir, entryPointScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fe.WriteString("#!/bin/sh\n" + st.task.BuildOptions.EntrypointCmd)
	if err != nil {
		return err
	}
	defer fe.Close()
	defer os.Remove(fe.Name())

	var tmp *os.File
	tmp, err = os.CreateTemp("", "Dockerfile")
	if err != nil {
		return err
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())
	dockerfile := fmt.Sprintf(
		`FROM %s
WORKDIR /srv
COPY . .
RUN ["/srv/%s"]
ENTRYPOINT ["/srv/%s"]`,
		st.task.BuildOptions.BaseImageName,
		buildScriptName,
		entryPointScriptName,
	)
	if _, err = tmp.WriteString(dockerfile); err != nil {
		return err
	}

	_, err = s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:  buildkit.ExporterImage,
			Attrs: exportAttrs,
		}},
		LocalDirs: map[string]string{
			"context":    st.repositoryTempDir,
			"dockerfile": filepath.Dir(tmp.Name()),
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": filepath.Base(tmp.Name())},
		Session:       []session.Attachable{authprovider.NewDockerAuthProvider(io.Discard)},
	}, ch)
	return err
}

func (s *builderService) buildStatic(ctx context.Context, st *state) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if len(st.task.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
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
	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return st.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"context": st.repositoryTempDir,
		},
	}, ch)
	return err
}

func (s *builderService) buildStaticWithConfig(
	ctx context.Context,
	st *state,
	ch chan *buildkit.SolveStatus,
) error {
	b := llb.Image(st.task.BuildOptions.BaseImageName).
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			AllowWildcard:  true,
			CreateDestPath: true,
		})).
		Run(llb.Shlex(st.task.BuildOptions.BuildCmd)).
		Root()
	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(b, st.task.BuildOptions.ArtifactPath, "/", &llb.CopyInfo{
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

type state struct {
	task       *builder.Task
	repository *domain.Repository

	repositoryTempDir string
	logTempFile       *os.File
	artifactTempFile  *os.File
}

func newState(task *builder.Task, repo *domain.Repository) *state {
	return &state{
		task:       task,
		repository: repo,
	}
}

func (s *state) initTempFiles(useArtifactTempFile bool) error {
	var err error

	// ログ用一時ファイル作成
	s.logTempFile, err = os.CreateTemp("", "buildlog-")
	if err != nil {
		return errors.Wrap(err, "failed to create tmp log file")
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
	if s.logTempFile == nil {
		return io.Discard
	}
	return s.logTempFile
}

func (s *state) writeLog(a ...interface{}) {
	_, _ = fmt.Fprintln(s.getLogWriter(), a...)
}

func toPBSettleReason(status builder.BuildStatus) pb.BuildSettled_Reason {
	switch status {
	case builder.BuildStatusSucceeded:
		return pb.BuildSettled_SUCCESS
	case builder.BuildStatusFailed:
		return pb.BuildSettled_FAILED
	case builder.BuildStatusCanceled:
		return pb.BuildSettled_CANCELLED
	default:
		panic(fmt.Sprintf("unexpected settled status: %v", status))
	}
}
