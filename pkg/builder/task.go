package builder

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/leandro-lugaresi/hub"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/idgen"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Task struct {
	Static        bool
	BuildID       string
	RepositoryURL string
	ImageName     string
	BuildLogM     models.BuildLog

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
	_, err = git.PlainCloneContext(ctx, t.repositoryTempDir, false, &git.CloneOptions{URL: t.RepositoryURL, Depth: 1})
	if err != nil {
		_ = os.RemoveAll(t.repositoryTempDir)
		log.WithError(err).Errorf("failed to clone repository: %s", t.RepositoryURL)
		return err
	}

	// ビルドログのエントリをDBに挿入
	t.BuildLogM.ID = t.BuildID
	if err := t.BuildLogM.Insert(ctx, s.db, boil.Infer()); err != nil {
		log.WithError(err).Errorf("failed to insert build_log entry (buildID: %d)", t.BuildID)
		return err
	}

	// 実行
	t.ctx, t.cancelFunc = context.WithCancel(context.Background())
	go s.processTask(t)
	s.bus.Publish(hub.Message{
		Name: IEventBuildStarted,
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
	// TODO ログファイルの保存実装をStorageインターフェースで抽象化する
	if err := util.MoveFile(t.logTempFile.Name(), filepath.Join("/neoshowcase/buildlogs", t.BuildID)); err != nil {
		log.WithError(err).Errorf("failed to save build log (%s)", t.BuildID)
	}

	if t.Static {
		_ = t.artifactTempFile.Close()
		if result == models.BuildLogsResultSUCCEEDED {
			sid := idgen.New()

			// 生成物tarの保存
			stat, _ := os.Stat(t.artifactTempFile.Name())
			artifact := models.Artifact{
				ID:         sid,
				BuildLogID: t.BuildID,
				Size:       stat.Size(),
				CreatedAt:  time.Now(),
			}

			// TODO 生成物tarの保存をStorageインターフェースで抽象化する
			if err := util.MoveFile(t.artifactTempFile.Name(), filepath.Join("/neoshowcase/artifacts", fmt.Sprintf("%s.tar", sid))); err != nil {
				log.WithError(err).Errorf("failed to save artifact tar file (BuildID: %s, ArtifactID: %s)", t.BuildID, sid)
			}

			if err := artifact.Insert(context.Background(), s.db, boil.Infer()); err != nil {
				log.WithError(err).Errorf("failed to insert artifact entry (BuildID: %s, ArtifactID: %s)", t.BuildID, sid)
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
			Name: IEventBuildFailed,
			Fields: hub.Fields{
				"task": t,
			},
		})
	case models.BuildLogsResultCANCELED:
		s.bus.Publish(hub.Message{
			Name: IEventBuildCanceled,
			Fields: hub.Fields{
				"task": t,
			},
		})
	case models.BuildLogsResultSUCCEEDED:
		s.bus.Publish(hub.Message{
			Name: IEventBuildSucceeded,
			Fields: hub.Fields{
				"task": t,
			},
		})
	default:
		panic(result)
	}

	return nil
}
