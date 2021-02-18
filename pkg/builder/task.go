package builder

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/event"
	"golang.org/x/sync/errgroup"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/traPtitech/neoshowcase/pkg/storage"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	startupScriptName    = "shell.sh"
	entryPointScriptName = "entrypoint.sh"
)

type Task struct {
	Static       bool
	BuildID      string
	BuildSource  *api.BuildSource
	BuildOptions *api.BuildOptions
	ImageName    string
	BuildLogM    models.BuildLog

	ctx               context.Context
	cancelFunc        func()
	repositoryTempDir string
	logTempFile       *os.File
	artifactTempFile  *os.File
}

func (t *Task) buildLogWriter() io.Writer {
	return t.logTempFile
}

func (t *Task) artifactWriter() io.WriteCloser {
	return t.artifactTempFile
}

func (t *Task) writeLog(a ...interface{}) {
	_, _ = fmt.Fprintln(t.buildLogWriter(), a...)
}

func (t *Task) startAsync(ctx context.Context, s *Service) error {
	// ログ用一時ファイル作成
	logF, err := ioutil.TempFile("", "buildlog")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary log file")
		return err
	}
	t.logTempFile = logF

	// 成果物tarの一時保存先作成
	if t.Static {
		artF, err := ioutil.TempFile("", "artifacts")
		if err != nil {
			log.WithError(err).Errorf("failed to create temporary artifact file")
			return err
		}
		t.artifactTempFile = artF
	}

	// リポジトリクローン用の一時ディレクトリ作成
	dir, err := ioutil.TempDir("", "repo")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary repository directory")
		return err
	}
	t.repositoryTempDir = dir

	// リポジトリをクローン
	refName := plumbing.HEAD
	if t.BuildSource.Ref != "" {
		refName = plumbing.ReferenceName("refs/" + t.BuildSource.Ref)
	}
	_, err = git.PlainCloneContext(ctx, t.repositoryTempDir, false, &git.CloneOptions{URL: t.BuildSource.RepositoryUrl, Depth: 1, ReferenceName: refName})
	if err != nil {
		_ = os.RemoveAll(t.repositoryTempDir)
		log.WithError(err).Errorf("failed to clone repository: %s", t.BuildSource.RepositoryUrl)
		return err
	}

	// ビルドログのエントリをDBに挿入
	t.BuildLogM.ID = t.BuildID
	if err := t.BuildLogM.Insert(ctx, s.db, boil.Infer()); err != nil {
		log.WithError(err).Errorf("failed to insert build_log entry (buildID: %s)", t.BuildID)
		return err
	}

	// 実行
	t.ctx, t.cancelFunc = context.WithCancel(context.Background())
	go s.processTask(t)
	s.bus.Publish(hub.Message{
		Name: event.BuilderBuildStarted,
		Fields: hub.Fields{
			"task": t,
		},
	})
	return nil
}

func (t *Task) postProcess(s *Service, result string) error {
	log.WithField("buildID", t.BuildID).
		WithField("result", result).
		Debugf("task finished")
	t.cancelFunc()

	// ログファイルの保存
	_ = t.logTempFile.Close()
	if err := storage.SaveLogFile(s.storage, t.logTempFile.Name(), filepath.Join("buildlogs", t.BuildID), t.BuildID); err != nil {
		log.WithError(err).Errorf("failed to save build log (%s)", t.BuildID)
	}

	if t.Static {
		// 生成物tarの保存
		_ = t.artifactTempFile.Close()
		if result == models.BuildLogsResultSUCCEEDED {
			sid := idgen.New()
			err := storage.SaveArtifact(s.storage, t.artifactTempFile.Name(), filepath.Join("artifacts", fmt.Sprintf("%s.tar", sid)), s.db, t.BuildID, sid)
			if err != nil {
				log.WithError(err).Errorf("failed to save directory to tar (BuildID: %s, ArtifactID: %s)", t.BuildID, sid)
			}
		} else {
			_ = os.Remove(t.artifactTempFile.Name())
		}
	}

	// 一時リポジトリディレクトリの削除
	_ = os.RemoveAll(t.repositoryTempDir)

	// BuildLog更新
	t.BuildLogM.Result = result
	t.BuildLogM.FinishedAt = null.TimeFrom(time.Now())
	if _, err := t.BuildLogM.Update(context.Background(), s.db, boil.Infer()); err != nil {
		log.WithError(err).Errorf("failed to update build_log entry (%s)", t.BuildID)
	}

	// イベント発行
	switch result {
	case models.BuildLogsResultFAILED:
		s.bus.Publish(hub.Message{
			Name: event.BuilderBuildFailed,
			Fields: hub.Fields{
				"task": t,
			},
		})
	case models.BuildLogsResultCANCELED:
		s.bus.Publish(hub.Message{
			Name: event.BuilderBuildCanceled,
			Fields: hub.Fields{
				"task": t,
			},
		})
	case models.BuildLogsResultSUCCEEDED:
		s.bus.Publish(hub.Message{
			Name: event.BuilderBuildSucceeded,
			Fields: hub.Fields{
				"task": t,
			},
		})
	default:
		panic(result)
	}

	return nil
}

func (t *Task) buildImage(s *Service) error {
	logWriter := t.buildLogWriter()
	if logWriter == nil {
		logWriter = ioutil.Discard // ログを破棄
	}

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(t.ctx)
	eg.Go(func() (err error) {
		// イメージの出力先設定
		exportAttrs := map[string]string{}
		if len(t.ImageName) == 0 {
			// ImageNameの指定がない場合はビルドするだけで、イメージを保存しない
			exportAttrs["name"] = "build-" + t.BuildID
		} else {
			exportAttrs["name"] = s.config.Buildkit.Registry + "/" + t.ImageName + ":" + t.BuildID
			exportAttrs["push"] = "true"
		}

		if t.BuildOptions == nil || len(t.BuildOptions.BaseImageName) == 0 {
			// リポジトリルートのDockerfileを使用
			// entrypoint, startupコマンドは無視
			_, err = s.buildkit.Solve(ctx, nil, client.SolveOpt{
				Exports: []client.ExportEntry{{
					Type:  client.ExporterImage,
					Attrs: exportAttrs,
				}},
				LocalDirs: map[string]string{
					builder.DefaultLocalNameContext:    t.repositoryTempDir,
					builder.DefaultLocalNameDockerfile: t.repositoryTempDir,
				},
				Frontend:      "dockerfile.v0",
				FrontendAttrs: map[string]string{"filename": "Dockerfile"},
				Session:       []session.Attachable{authprovider.NewDockerAuthProvider(ioutil.Discard)},
			}, ch)
		} else {
			// 指定したベースイメージを使用
			var fs, fe *os.File
			fs, err := ioutil.TempFile("", startupScriptName)
			if err != nil {
				return err
			}
			cmd := fmt.Sprintf(`
#!/bin/sh

%s
`, t.BuildOptions.StartupCmd)
			_, err = fs.WriteString(cmd)
			if err != nil {
				return err
			}
			defer fs.Close()
			defer os.Remove(fs.Name())

			fe, err = ioutil.TempFile("", entryPointScriptName)
			if err != nil {
				return err
			}
			cmd = fmt.Sprintf(`
#!/bin/sh

%s
`, t.BuildOptions.EntrypointCmd)
			_, err = fe.WriteString(cmd)
			if err != nil {
				return err
			}
			defer fe.Close()
			defer os.Remove(fe.Name())

			dockerfile := fmt.Sprintf(`
FROM %s
COPY . .
RUN ./startup.sh
ENTRYPOINT ./entrypoint.sh
`, t.BuildOptions.BaseImageName)

			var tmp *os.File
			tmp, err = ioutil.TempFile("", "Dockerfile")
			if err != nil {
				return err
			}
			defer tmp.Close()
			defer os.Remove(tmp.Name())
			if _, err := tmp.WriteString(dockerfile); err != nil {
				return err
			}

			_, err = s.buildkit.Solve(ctx, nil, client.SolveOpt{
				Exports: []client.ExportEntry{{
					Type:  client.ExporterImage,
					Attrs: exportAttrs,
				}},
				LocalDirs: map[string]string{
					builder.DefaultLocalNameContext:    t.repositoryTempDir,
					builder.DefaultLocalNameDockerfile: filepath.Dir(tmp.Name()),
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
		return progressui.DisplaySolveStatus(ctx, "", nil, logWriter, ch)
	})

	return eg.Wait()
}

func (t *Task) buildStatic(s *Service) error {
	logWriter := t.buildLogWriter()
	if logWriter == nil {
		logWriter = ioutil.Discard // ログを破棄
	}

	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(t.ctx)
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

			_, err = s.buildkit.Solve(ctx, def, client.SolveOpt{
				Exports: []client.ExportEntry{{
					Type:   client.ExporterTar,
					Output: func(_ map[string]string) (io.WriteCloser, error) { return t.artifactWriter(), nil },
				}},
				LocalDirs: map[string]string{
					"local-src": t.repositoryTempDir,
				},
			}, ch)
		}
		return err
	})
	eg.Go(func() error {
		// ビルドログを収集
		return progressui.DisplaySolveStatus(ctx, "", nil, logWriter, ch)
	})

	return eg.Wait()
}
