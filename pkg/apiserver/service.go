package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/leandro-lugaresi/hub"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/container/dockerimpl"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Service struct {
	server *grpc.Server
	db     *sql.DB
	bus    *hub.Hub

	builderConn      *grpc.ClientConn
	ssgenConn        *grpc.ClientConn
	containerManager container.Manager

	config Config
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
		bus:    hub.New(),
	}

	reflection.Register(s.server)

	// DBに接続
	db, err := c.DB.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}
	s.db = db

	// Builderに接続
	builderConn, err := c.Builder.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect builder service: %w", err)
	}
	s.builderConn = builderConn

	// SSGenに接続
	ssgenConn, err := c.SSGen.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect ssgen service: %w", err)
	}
	s.ssgenConn = ssgenConn

	// コンテナマネージャー生成
	connM, err := dockerimpl.NewManager(s.bus)
	if err != nil {
		return nil, fmt.Errorf("failed to init container manager: %w", err)
	}
	s.containerManager = connM

	return s, nil
}

func (s *Service) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPC.GetPort()))
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	return s.server.Serve(listener)
}

func (s *Service) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s.server.GracefulStop()
		return nil
	})
	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.builderConn.Close()
	})
	eg.Go(func() error {
		return s.ssgenConn.Close()
	})
	eg.Go(func() error {
		return s.containerManager.Dispose(ctx)
	})

	return eg.Wait()
}

func (s *Service) builder() builderApi.BuilderServiceClient {
	return builderApi.NewBuilderServiceClient(s.builderConn)
}

func (s *Service) ssgen() ssgenApi.StaticSiteGenServiceClient {
	return ssgenApi.NewStaticSiteGenServiceClient(s.ssgenConn)
}
