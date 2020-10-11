package staticsitegen

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	server *grpc.Server
	engine generator.Engine
	db     *sql.DB

	config Config
	api.UnimplementedStaticSiteGenServiceServer
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
	}
	api.RegisterStaticSiteGenServiceServer(s.server, s)

	engine, err := c.GetEngine()
	if err != nil {
		return nil, err
	}
	s.engine = engine

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

	return eg.Wait()
}
