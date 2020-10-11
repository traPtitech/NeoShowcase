package builder

import (
	"context"
	"database/sql"
	"fmt"
	buildkit "github.com/moby/buildkit/client"
	"github.com/traPtitech/neoshowcase/pkg/builder/api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Service struct {
	server *grpc.Server
	client *buildkit.Client
	db     *sql.DB

	config Config
	api.UnimplementedBuilderServiceServer
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
	}
	api.RegisterBuilderServiceServer(s.server, s)

	// buildkitdに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, c.Buildkit.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Buildkit Client: %w", err)
	}
	s.client = client

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
	eg.Go(func() error {
		return s.client.Close()
	})

	return eg.Wait()
}
