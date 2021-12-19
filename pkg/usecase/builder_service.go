package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/status"
)

const (
	startupScriptName    = "shell.sh"
	entryPointScriptName = "entrypoint.sh"
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

	// TODO 後で消す
	db *sql.DB

	registry string

	task              *builder.Task
	internalTaskState *internalTaskState
	taskLock          sync.Mutex

	status     builder.State
	statusLock sync.RWMutex
}

func NewBuilderService(buildkit *buildkit.Client, storage domain.Storage, eventbus domain.Bus, db *sql.DB, registry builder.DockerImageRegistryString) BuilderService {
	return &builderService{
		buildkit: buildkit,
		storage:  storage,
		eventbus: eventbus,
		db:       db,
		registry: string(registry),

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

func (s *builderService) CancelBuild(ctx context.Context) (string, error) {
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

func (s *builderService) Shutdown(ctx context.Context) error {
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
		BuildLogM: models.BuildLog{
			ID:        task.BuildID,
			Result:    builder.BuildStatusBuilding,
			StartedAt: time.Now(),
			BranchID:  task.BranchID.String,
		},
	}

	// ログ用一時ファイル作成
	logF, err := ioutil.TempFile("", "buildlog")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary log file")
		return fmt.Errorf("failed to create tmp log file: %w", err)
	}
	intState.logTempFile = logF

	// 成果物tarの一時保存先作成
	if task.Static {
		artF, err := ioutil.TempFile("", "artifacts")
		if err != nil {
			log.WithError(err).Errorf("failed to create temporary artifact file")
			return fmt.Errorf("failed to create tmp artifact file: %w", err)
		}
		intState.artifactTempFile = artF
	}

	// リポジトリクローン用の一時ディレクトリ作成
	dir, err := ioutil.TempDir("", "repo")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary repository directory")
		return fmt.Errorf("failed to create tmp repository dir: %w", err)
	}
	intState.repositoryTempDir = dir

	// リポジトリをクローン
	refName := plumbing.HEAD
	if task.BuildSource.Ref != "" {
		refName = plumbing.ReferenceName("refs/" + task.BuildSource.Ref)
	}
	_, err = git.PlainCloneContext(ctx, intState.repositoryTempDir, false, &git.CloneOptions{URL: task.BuildSource.RepositoryUrl, Depth: 1, ReferenceName: refName})
	if err != nil {
		_ = os.RemoveAll(intState.repositoryTempDir)
		log.WithError(err).Errorf("failed to clone repository: %s", task.BuildSource.RepositoryUrl)
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// TODO: QueueBuildに移す
	// ビルドログのエントリをDBに挿入
	if err := intState.BuildLogM.Insert(ctx, s.db, boil.Infer()); err != nil {
		log.WithError(err).Errorf("failed to insert build_log entry (buildID: %s)", task.BuildID)
		return fmt.Errorf("failed to save BuildLog: %w", err)
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

	result := builder.BuildStatusFailed
	// 後処理関数
	defer func() {
		// タスク破棄
		s.taskLock.Lock()
		s.task = nil
		s.internalTaskState = nil
		s.taskLock.Unlock()

		log.WithField("buildID", task.BuildID).
			WithField("result", result).
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
			if result == builder.BuildStatusSucceeded {
				sid := domain.NewID()
				err := domain.SaveArtifact(s.storage, intState.artifactTempFile.Name(), filepath.Join("artifacts", fmt.Sprintf("%s.tar", sid)), s.db, task.BuildID, sid)
				if err != nil {
					log.WithError(err).Errorf("failed to save directory to tar (BuildID: %s, ArtifactID: %s)", task.BuildID, sid)
				}
			} else {
				_ = os.Remove(intState.artifactTempFile.Name())
			}
		}

		// 一時リポジトリディレクトリの削除
		_ = os.RemoveAll(intState.repositoryTempDir)

		// BuildLog更新
		intState.BuildLogM.Result = result
		intState.BuildLogM.FinishedAt = null.TimeFrom(time.Now())
		if _, err := intState.BuildLogM.Update(context.Background(), s.db, boil.Infer()); err != nil {
			log.WithError(err).Errorf("failed to update build_log entry (%s)", task.BuildID)
		}

		// イベント発行
		switch result {
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
		if err == context.Canceled || err == context.DeadlineExceeded || errors.Is(err, status.FromContextError(context.Canceled).Err()) {
			result = builder.BuildStatusCanceled
			intState.writeLog("CANCELED")
			return
		}
		result = builder.BuildStatusFailed
		return
	}

	// 成功
	intState.writeLog("BUILD SUCCESSFUL")
	result = builder.BuildStatusSucceeded
}

func (s *builderService) buildImage(t *builder.Task, intState *internalTaskState) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(intState.ctx)
	eg.Go(func() (err error) {
		// イメージの出力先設定
		exportAttrs := map[string]string{}
		if len(t.ImageName) == 0 {
			// ImageNameの指定がない場合はビルドするだけで、イメージを保存しない
			exportAttrs["name"] = "build-" + t.BuildID
		} else {
			exportAttrs["name"] = s.registry + "/" + t.ImageName + ":" + t.BuildID
			exportAttrs["push"] = "true"
		}

		if t.BuildOptions == nil || len(t.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
			// entrypoint, startupコマンドは無視
			_, err = s.buildkit.Solve(ctx, nil, buildkit.SolveOpt{
				Exports: []buildkit.ExportEntry{{
					Type:  buildkit.ExporterImage,
					Attrs: exportAttrs,
				}},
				LocalDirs: map[string]string{
					"context":    intState.repositoryTempDir,
					"dockerfile": intState.repositoryTempDir,
				},
				Frontend:      "dockerfile.v0",
				FrontendAttrs: map[string]string{"filename": "Dockerfile"},
				Session:       []session.Attachable{authprovider.NewDockerAuthProvider(ioutil.Discard)},
			}, ch)
		} else {
			// 指定したベースイメージを使用
			var fs, fe *os.File
			fs, err = os.OpenFile(filepath.Join(intState.repositoryTempDir, startupScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			}
			scmd := fmt.Sprintf(`#!/bin/sh
%s
`, t.BuildOptions.StartupCmd)
			_, err = fs.WriteString(scmd)
			if err != nil {
				return err
			}
			defer fs.Close()
			defer os.Remove(fs.Name())

			fe, err = os.OpenFile(filepath.Join(intState.repositoryTempDir, entryPointScriptName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			}
			ecmd := fmt.Sprintf(`#!/bin/sh
%s
`, t.BuildOptions.EntrypointCmd)
			_, err = fe.WriteString(ecmd)
			if err != nil {
				return err
			}
			defer fe.Close()
			defer os.Remove(fe.Name())

			dockerfile := fmt.Sprintf(`
FROM %s
COPY . .
RUN ./%s
ENTRYPOINT ./%s
`, t.BuildOptions.BaseImageName, startupScriptName, entryPointScriptName)
			var tmp *os.File
			tmp, err = ioutil.TempFile("", "Dockerfile")
			if err != nil {
				return err
			}
			defer tmp.Close()
			defer os.Remove(tmp.Name())
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
				Session:       []session.Attachable{authprovider.NewDockerAuthProvider(ioutil.Discard)},
			}, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(context.TODO(), "", nil, intState.getLogWriter(), ch)
	})
	return eg.Wait()
}

func (s *builderService) buildStatic(t *builder.Task, intState *internalTaskState) error {
	ch := make(chan *buildkit.SolveStatus)
	eg, ctx := errgroup.WithContext(intState.ctx)
	eg.Go(func() (err error) {
		if t.BuildOptions == nil || len(t.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
			// entrypoint, startupコマンドは無視
			// TODO
			panic("not implemented")
		} else {
			// 指定したベースイメージを使用
			b := llb.Image(t.BuildOptions.BaseImageName).
				File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
					AllowWildcard:  true,
					CreateDestPath: true,
				})).
				Run(llb.Shlex(t.BuildOptions.StartupCmd)).
				Root()
			// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
			def, _ := llb.
				Scratch().
				File(llb.Copy(b, t.BuildOptions.ArtifactPath, "/", &llb.CopyInfo{
					CopyDirContentsOnly: true,
					CreateDestPath:      true,
					AllowWildcard:       true,
				})).
				Marshal(context.Background())

			_, err = s.buildkit.Solve(ctx, def, buildkit.SolveOpt{
				Exports: []buildkit.ExportEntry{{
					Type:   buildkit.ExporterTar,
					Output: func(_ map[string]string) (io.WriteCloser, error) { return intState.artifactTempFile, nil },
				}},
				LocalDirs: map[string]string{
					"local-src": intState.repositoryTempDir,
				},
			}, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(context.TODO(), "", nil, intState.getLogWriter(), ch)
	})
	return eg.Wait()
}

type internalTaskState struct {
	BuildLogM         models.BuildLog
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
