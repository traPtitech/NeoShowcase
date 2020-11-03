package builder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/leandro-lugaresi/hub"
	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"
)

type Service struct {
	server   *grpc.Server
	buildkit BuildkitWrapper
	db       *sql.DB
	bus      *hub.Hub

	config Config
	api.UnimplementedBuilderServiceServer

	task     *Task
	taskLock sync.Mutex

	state     api.BuilderStatus
	stateLock sync.RWMutex
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
		bus:    hub.New(),
	}
	api.RegisterBuilderServiceServer(s.server, s)
	reflection.Register(s.server)

	// buildkitdに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, c.Buildkit.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Buildkit Client: %w", err)
	}
	s.buildkit.client = client

	// DBに接続
	db, err := c.DB.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	s.db = db

	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.GetPort()))
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	s.stateLock.Lock()
	s.state = api.BuilderStatus_WAITING
	s.stateLock.Unlock()
	return s.server.Serve(listener)
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.stateLock.Lock()
	s.state = api.BuilderStatus_UNAVAILABLE
	s.stateLock.Unlock()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s.server.GracefulStop()
		return nil
	})
	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.buildkit.client.Close()
	})

	return eg.Wait()
}

func (s *Service) processTask(t *Task) {
	s.taskLock.Lock()
	s.task = t
	s.taskLock.Unlock()

	var f *os.File
	result := models.BuildLogsResultFAILED
	// 後処理関数
	defer func() {
		log.WithField("buildID", t.BuildID).
			WithField("result", result).
			Debugf("task finished")

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

		// タスク破棄
		s.taskLock.Lock()
		s.task = nil
		s.taskLock.Unlock()
		t.dispose()
		s.stateLock.Lock()
		s.state = api.BuilderStatus_WAITING
		s.stateLock.Unlock()

		// BuildLog更新
		t.BuildLogM.Result = result
		t.BuildLogM.FinishedAt = null.TimeFrom(time.Now())
		if _, err := t.BuildLogM.Update(context.Background(), s.db, boil.Infer()); err != nil {
			log.WithError(err).Errorf("failed to update build_log entry (%d)", t.BuildID)
		}

		// ログファイルの保存 TODO
		if f != nil {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}
	}()

	// ログ用一時ファイル作成
	f, err := ioutil.TempFile("", "buildlog")
	if err != nil {
		log.WithError(err).Errorf("failed to create temporary log file (buildID: %d)", t.BuildID)
		return
	}

	// リポジトリをクローン
	dir, err := ioutil.TempDir("", "repo")
	if err != nil {
		result = models.BuildLogsResultFAILED
		log.WithError(err).Errorf("failed to create temporary repository directory (buildID: %d)", t.BuildID)
		_, _ = fmt.Fprintln(f, "[INTERNAL ERROR OCCURRED]")
		return
	}
	defer os.RemoveAll(dir)

	_, _ = fmt.Fprint(f, "CLONE REPOSITORY...")
	_, err = git.PlainCloneContext(t.Ctx, dir, false, &git.CloneOptions{URL: t.RepositoryURL, Depth: 1})
	if err != nil {
		log.Debug(err)
		if err == context.Canceled || err == context.DeadlineExceeded {
			result = models.BuildLogsResultCANCELED
			_, _ = fmt.Fprintln(f, "CANCELED")
			return
		}
		result = models.BuildLogsResultFAILED
		_, _ = fmt.Fprintln(f, "ERROR")
		_, _ = fmt.Fprintln(f, err.Error())
		return
	}
	_, _ = fmt.Fprintln(f, "SUCCESS")

	// ビルド
	_, _ = fmt.Fprintln(f, "START BUILDING")
	if err := s.buildkit.BuildImage(t.Ctx, BuildImageArgs{
		ImageName:  s.config.Buildkit.Registry + "/" + t.ImageName,
		ContextDir: dir,
	}, f); err != nil {
		log.Debug(err)
		if err == context.Canceled || err == context.DeadlineExceeded {
			result = models.BuildLogsResultCANCELED
			_, _ = fmt.Fprintln(f, "CANCELED")
			return
		}
		_, _ = fmt.Fprintln(f, err.Error())
		return
	}

	// 成功
	_, _ = fmt.Fprintln(f, "BUILD SUCCESSFUL")
	result = models.BuildLogsResultSUCCEEDED
	return
}
