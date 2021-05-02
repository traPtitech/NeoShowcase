package builder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/leandro-lugaresi/hub"
	buildkit "github.com/moby/buildkit/client"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	storage2 "github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Service struct {
	server   *grpc.Server
	buildkit *buildkit.Client
	db       *sql.DB
	bus      *hub.Hub
	storage  storage2.Storage

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
	s.buildkit = client

	// DBに接続
	db, err := c.DB.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	s.db = db

	// Storageに接続
	storage, err := c.Storage.InitStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	s.storage = storage

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
		return s.buildkit.Close()
	})

	return eg.Wait()
}

func (s *Service) processTask(t *Task) {
	s.setTask(t)

	result := models.BuildLogsResultFAILED
	// 後処理関数
	defer func() {
		// タスク破棄
		s.setTask(nil)

		t.postProcess(s, result)

		s.stateLock.Lock()
		s.state = api.BuilderStatus_WAITING
		s.stateLock.Unlock()
	}()

	// ビルド
	t.writeLog("START BUILDING")
	var err error
	if t.Static {
		// 静的ファイルビルド
		err = t.buildStatic(s)
	} else {
		// DockerImageビルド
		err = t.buildImage(s)
	}
	if err != nil {
		log.Debug(err)
		if err == context.Canceled || err == context.DeadlineExceeded || errors.Is(err, status.FromContextError(context.Canceled).Err()) {
			result = models.BuildLogsResultCANCELED
			t.writeLog("CANCELED")
			return
		}
		result = models.BuildLogsResultFAILED
		return
	}

	// 成功
	t.writeLog("BUILD SUCCESSFUL")
	result = models.BuildLogsResultSUCCEEDED
	return
}

func (s *Service) setTask(t *Task) {
	s.taskLock.Lock()
	s.task = t
	s.taskLock.Unlock()
}
