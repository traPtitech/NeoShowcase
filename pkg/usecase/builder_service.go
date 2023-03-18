package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
)

const (
	buildScriptName      = "neoshowcase_internal_build.sh"
	entryPointScriptName = "neoshowcase_internal_entrypoint.sh"
)

type BuilderService interface {
	GetStatus() builder.State
	StreamEvents() domain.Subscription
	StartBuild(ctx context.Context, task *builder.Task) error
	CancelBuild(ctx context.Context) (buildID string, err error)
	Shutdown(ctx context.Context) error
}

type builderService struct {
	buildkit *buildkit.Client
	storage  domain.Storage
	eventbus domain.Bus

	// TODO: 後で消す
	artifactRepo repository.ArtifactRepository
	buildRepo    repository.BuildRepository

	task              *builder.Task
	internalTaskState *internalTaskState
	taskLock          sync.Mutex

	status     builder.State
	statusLock sync.RWMutex
}

func NewBuilderService(buildkit *buildkit.Client, storage domain.Storage, eventbus domain.Bus, artifactRepo repository.ArtifactRepository, buildRepo repository.BuildRepository) BuilderService {
	return &builderService{
		buildkit:     buildkit,
		storage:      storage,
		eventbus:     eventbus,
		artifactRepo: artifactRepo,
		buildRepo:    buildRepo,

		status: builder.StateWaiting,
	}
}

func (s *builderService) GetStatus() builder.State {
	s.statusLock.RLock()
	defer s.statusLock.RUnlock()
	return s.status
}

func (s *builderService) StreamEvents() domain.Subscription {
	return s.eventbus.Subscribe(event.BuilderBuildStarted, event.BuilderBuildFailed, event.BuilderBuildCanceled, event.BuilderBuildSucceeded)
}

func (s *builderService) StartBuild(ctx context.Context, task *builder.Task) error {
	s.statusLock.Lock()
	if s.status != builder.StateWaiting {
		s.statusLock.Unlock()
		return fmt.Errorf("builder unavailable")
	}

	if err := s.initializeTask(ctx, task); err != nil {
		return fmt.Errorf("failed to initialize Task: %w", err)
	}

	s.status = builder.StateBuilding
	s.statusLock.Unlock()
	return nil
}

func (s *builderService) CancelBuild(_ context.Context) (string, error) {
	s.statusLock.RLock()
	state := s.status
	s.statusLock.RUnlock()

	if state != builder.StateBuilding {
		return "", nil
	}

	s.taskLock.Lock()
	defer s.taskLock.Unlock()
	buildID := s.task.BuildID
	s.internalTaskState.cancelFunc()
	return buildID, nil
}

func (s *builderService) Shutdown(_ context.Context) error {
	s.statusLock.Lock()
	prevStatus := s.status
	s.status = builder.StateUnavailable
	s.statusLock.Unlock()

	if prevStatus != builder.StateBuilding {
		return nil
	}

	s.taskLock.Lock()
	intState := s.internalTaskState
	s.taskLock.Unlock()

	<-intState.ctx.Done()
	return nil
}

func (s *builderService) initializeTask(ctx context.Context, task *builder.Task) error {
	intState := &internalTaskState{
		Build: domain.Build{
			ID:            task.BuildID,
			Status:        builder.BuildStatusBuilding,
			ApplicationID: task.ApplicationID,
		},
	}

	// ログ用一時ファイル作成
	logF, err := os.CreateTemp("", "buildlog")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary log file")
		return fmt.Errorf("failed to create tmp log file: %w", err)
	}
	intState.logTempFile = logF

	// 成果物tarの一時保存先作成
	if task.Static {
		artF, err := os.CreateTemp("", "artifacts")
		if err != nil {
			log.WithError(err).Errorf("failed to create temporary artifact file")
			return fmt.Errorf("failed to create tmp artifact file: %w", err)
		}
		intState.artifactTempFile = artF
	}

	// リポジトリクローン用の一時ディレクトリ作成
	dir, err := os.MkdirTemp("", "repo")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary repository directory")
		return fmt.Errorf("failed to create tmp repository dir: %w", err)
	}
	intState.repositoryTempDir = dir

	// リポジトリをクローン
	repo, err := git.PlainInit(intState.repositoryTempDir, false)
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to init repository: %s", intState.repositoryTempDir)
		return fmt.Errorf("failed to init repository: %w", err)
	}
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{task.BuildSource.RepositoryUrl},
	})
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to add remote: %s", task.BuildSource.RepositoryUrl)
		return fmt.Errorf("failed to add remote: %w", err)
	}
	targetRef := plumbing.NewRemoteReferenceName("origin", "target")
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+%s:%s", task.BuildSource.Commit, targetRef))},
		Depth:      1,
	})
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to clone repository: %s", task.BuildSource.RepositoryUrl)
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to get worktree: %s", intState.repositoryTempDir)
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	err = wt.Checkout(&git.CheckoutOptions{Branch: targetRef})
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to checkout: %s", intState.repositoryTempDir)
		return fmt.Errorf("failed to checkout: %w", err)
	}

	// Status を Building に変更
	args := repository.UpdateBuildArgs{
		ID:     intState.Build.ID,
		Status: intState.Build.Status,
	}
	if err := s.buildRepo.UpdateBuild(ctx, args); err != nil {
		log.WithError(err).Errorf("failed to update build_log entry (buildID: %s)", task.BuildID)
		return fmt.Errorf("failed to save Build: %w", err)
	}

	intState.ctx, intState.cancelFunc = context.WithCancel(context.Background())
	go s.processTask(task, intState)
	s.eventbus.Publish(event.BuilderBuildStarted, domain.Fields{
		"task": task,
	})
	return nil
}

func (s *builderService) processTask(task *builder.Task, intState *internalTaskState) {
	s.taskLock.Lock()
	s.task = task
	s.internalTaskState = intState
	s.taskLock.Unlock()

	status := builder.BuildStatusFailed
	// 後処理関数
	defer func() {
		// タスク破棄
		s.taskLock.Lock()
		s.task = nil
		s.internalTaskState = nil
		s.taskLock.Unlock()

		log.WithField("buildID", task.BuildID).
			WithField("status", status).
			Debugf("task finished")
		intState.cancelFunc()

		// ログファイルの保存
		_ = intState.logTempFile.Close()
		if err := domain.SaveLogFile(s.storage, intState.logTempFile.Name(), filepath.Join("buildlogs", task.BuildID), task.BuildID); err != nil {
			log.WithError(err).Errorf("failed to save build log (%s)", task.BuildID)
		}

		if task.Static {
			// 生成物tarの保存
			_ = intState.artifactTempFile.Close()
			if status == builder.BuildStatusSucceeded {
				sid := domain.NewID()
				filename := intState.artifactTempFile.Name()
				err := domain.SaveArtifact(s.storage, filename, filepath.Join("artifacts", fmt.Sprintf("%s.tar", sid)))
				if err != nil {
					log.WithError(err).Errorf("failed to save directory to tar (BuildID: %s, ArtifactID: %s)", task.BuildID, sid)
				}

				// TODO: エラー処理
				stat, _ := os.Stat(filename)
				err = s.artifactRepo.CreateArtifact(context.Background(), stat.Size(), task.BuildID, sid)
				if err != nil {
					log.WithError(err).Errorf("failed to create artifact (BuildID: %s, ArtifactID: %s)", task.BuildID, sid)
				}
			} else {
				_ = os.Remove(intState.artifactTempFile.Name())
			}
		}

		// 一時リポジトリディレクトリの削除
		_ = os.RemoveAll(intState.repositoryTempDir)

		// BuildLog更新
		intState.Build.Status = status
		args := repository.UpdateBuildArgs{
			ID:     intState.Build.ID,
			Status: intState.Build.Status,
		}
		if err := s.buildRepo.UpdateBuild(context.Background(), args); err != nil {
			log.WithError(err).Errorf("failed to update build_log entry (%s)", task.BuildID)
		}

		// イベント発行
		switch status {
		case builder.BuildStatusFailed:
			s.eventbus.Publish(event.BuilderBuildFailed, domain.Fields{
				"task": task,
			})
		case builder.BuildStatusCanceled:
			s.eventbus.Publish(event.BuilderBuildCanceled, domain.Fields{
				"task": task,
			})
		case builder.BuildStatusSucceeded:
			s.eventbus.Publish(event.BuilderBuildSucceeded, domain.Fields{
				"task": task,
			})
		}

		s.statusLock.Lock()
		s.status = builder.StateWaiting
		s.statusLock.Unlock()
	}()

	// ビルド
	intState.writeLog("START BUILDING")
	var err error
	if task.Static {
		// 静的ファイルビルド
		err = s.buildStatic(task, intState)
	} else {
		// DockerImageビルド
		err = s.buildImage(task, intState)
	}
	if err != nil {
		log.Debug(err)
		if err == context.Canceled || err == context.DeadlineExceeded || errors.Is(err, gstatus.FromContextError(context.Canceled).Err()) {
			status = builder.BuildStatusCanceled
			intState.writeLog("CANCELED")
			return
		}
		status = builder.BuildStatusFailed
		return
	}

	// 成功
	intState.writeLog("BUILD SUCCESSFUL")
	status = builder.BuildStatusSucceeded
}

func (s *builderService) buildImage(t *builder.Task, intState *internalTaskState) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(intState.ctx)
	eg.Go(func() (err error) {
		// イメージの出力先設定
		exportAttrs := map[string]string{
			"name": t.ImageName + ":" + t.ImageTag,
			"push": "true",
		}
		if len(t.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
			// entrypoint, startupコマンドは無視
			err = s.buildImageWithDockerfile(ctx, t, intState, exportAttrs, ch)
		} else {
			// 指定したベースイメージを使用
			err = s.buildImageWithConfig(ctx, t, intState, exportAttrs, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, intState.getLogWriter(), ch)
		return err
	})
	return eg.Wait()
}

func (s *builderService) buildImageWithDockerfile(
	ctx context.Context,
	t *builder.Task,
	intState *internalTaskState,
	exportAttrs map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	_, err := s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:  buildkit.ExporterImage,
			Attrs: exportAttrs,
		}},
		LocalDirs: map[string]string{
			"context":    intState.repositoryTempDir,
			"dockerfile": intState.repositoryTempDir,
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": t.BuildOptions.DockerfileName},
		Session:       []session.Attachable{authprovider.NewDockerAuthProvider(io.Discard)},
	}, ch)
	return err
}

func (s *builderService) buildImageWithConfig(
	ctx context.Context,
	t *builder.Task,
	intState *internalTaskState,
	exportAttrs map[string]string,
	ch chan *buildkit.SolveStatus,
) error {
	var fs, fe *os.File
	fs, err := os.OpenFile(filepath.Join(intState.repositoryTempDir, buildScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fs.WriteString("#!/bin/sh\n" + t.BuildOptions.BuildCmd)
	if err != nil {
		return err
	}
	defer fs.Close()
	defer os.Remove(fs.Name())

	fe, err = os.OpenFile(filepath.Join(intState.repositoryTempDir, entryPointScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fe.WriteString("#!/bin/sh\n" + t.BuildOptions.EntrypointCmd)
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
		t.BuildOptions.BaseImageName,
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
			"context":    intState.repositoryTempDir,
			"dockerfile": filepath.Dir(tmp.Name()),
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: map[string]string{"filename": filepath.Base(tmp.Name())},
		Session:       []session.Attachable{authprovider.NewDockerAuthProvider(io.Discard)},
	}, ch)
	return err
}

func (s *builderService) buildStatic(t *builder.Task, intState *internalTaskState) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(intState.ctx)
	eg.Go(func() (err error) {
		if len(t.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
			// entrypoint, startupコマンドは無視
			err = s.buildStaticWithDockerfile(ctx, t, intState, ch)
		} else {
			// 指定したベースイメージを使用
			err = s.buildStaticWithConfig(ctx, t, intState, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		// TODO: VertexWarningを使う (LLBのどのvertexに問題があったか)
		_, err := progressui.DisplaySolveStatus(ctx, "", nil, intState.getLogWriter(), ch)
		return err
	})
	return eg.Wait()
}

func (s *builderService) buildStaticWithDockerfile(
	ctx context.Context,
	t *builder.Task,
	intState *internalTaskState,
	ch chan *buildkit.SolveStatus,
) error {
	dockerfile, err := os.ReadFile(filepath.Join(intState.repositoryTempDir, t.BuildOptions.DockerfileName))
	if err != nil {
		return err
	}
	b, _, _, err := dockerfile2llb.Dockerfile2LLB(context.Background(), dockerfile, dockerfile2llb.ConvertOpt{})
	if err != nil {
		return err
	}
	def, err := llb.
		Scratch().
		File(llb.Copy(*b, t.BuildOptions.ArtifactPath, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(context.Background())
	if err != nil {
		return err
	}
	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return intState.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"context": intState.repositoryTempDir,
		},
	}, ch)
	return err
}

func (s *builderService) buildStaticWithConfig(
	ctx context.Context,
	t *builder.Task,
	intState *internalTaskState,
	ch chan *buildkit.SolveStatus,
) error {
	b := llb.Image(t.BuildOptions.BaseImageName).
		Dir("/srv").
		File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
			AllowWildcard:  true,
			CreateDestPath: true,
		})).
		Run(llb.Shlex(t.BuildOptions.BuildCmd)).
		Root()
	// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
	def, err := llb.
		Scratch().
		File(llb.Copy(b, t.BuildOptions.ArtifactPath, "/", &llb.CopyInfo{
			CopyDirContentsOnly: true,
			CreateDestPath:      true,
			AllowWildcard:       true,
		})).
		Marshal(context.Background())
	if err != nil {
		return err
	}

	_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
		Exports: []buildkit.ExportEntry{{
			Type:   buildkit.ExporterTar,
			Output: func(_ map[string]string) (io.WriteCloser, error) { return intState.artifactTempFile, nil },
		}},
		LocalDirs: map[string]string{
			"local-src": intState.repositoryTempDir,
		},
	}, ch)
	return err
}

type internalTaskState struct {
	Build             domain.Build
	ctx               context.Context
	cancelFunc        func()
	repositoryTempDir string
	logTempFile       *os.File
	artifactTempFile  *os.File
}

func (i *internalTaskState) getLogWriter() io.Writer {
	if i.logTempFile == nil {
		return io.Discard
	}
	return i.logTempFile
}

func (i *internalTaskState) writeLog(a ...interface{}) {
	_, _ = fmt.Fprintln(i.getLogWriter(), a...)
}
