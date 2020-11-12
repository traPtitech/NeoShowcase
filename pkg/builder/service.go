package builder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/models"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
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

		// TODO ビルドステージの構成を何らかの形で指定できるようにする
		builder := llb.Image("docker.io/library/node:14.11.0-alpine"). // FROM node:14.11.0-alpine as builder
			Dir("/app"). // WORKDIR /app
			File(llb.Copy(llb.Local("local-src"), "package*.json", "./", &llb.CopyInfo{
				AllowWildcard:  true,
				CreateDestPath: true,
			})). // COPY package*.json ./
			Run(llb.Shlex("npm i")). // RUN npm i
			File(llb.Copy(llb.Local("local-src"), ".", ".", &llb.CopyInfo{
				AllowWildcard:  true,
				CreateDestPath: true,
			})). // COPY . .
			Run(llb.Shlex("npm run build"), llb.AddEnv("NODE_ENV", "production")). // RUN NODE_ENV=production npm run build
			Root()

		// ビルドで生成された静的ファイルのみを含むScratchイメージを構成
		def, _ := llb.
			Scratch(). // FROM scratch
			File(llb.Copy(builder, "/app/dist", "/", &llb.CopyInfo{
				CopyDirContentsOnly: true,
				CreateDestPath:      true,
				AllowWildcard:       true,
			})). // COPY --from=builder /app/dist /
			Marshal(context.Background())

		// ビルド
		err = s.buildkit.BuildStatic(t.ctx, BuildStaticArgs{
			Output:     t.artifactWriter(),
			ContextDir: t.repositoryTempDir,
			LLB:        def,
		}, t.buildLogWriter())
	} else {
		// DockerImageビルド
		err = s.buildkit.BuildImage(t.ctx, BuildImageArgs{
			ImageName:  s.config.Buildkit.Registry + "/" + t.ImageName,
			ContextDir: t.repositoryTempDir,
		}, t.buildLogWriter())
	}
	if err != nil {
		log.Debug(err)
		if err == context.Canceled || err == context.DeadlineExceeded {
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