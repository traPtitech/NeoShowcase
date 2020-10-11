package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/traPtitech/neoshowcase/pkg/apiserver/api"
	builderApi "github.com/traPtitech/neoshowcase/pkg/builder/api"
	ssgenApi "github.com/traPtitech/neoshowcase/pkg/staticsitegen/api"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	server *grpc.Server
	db     *sql.DB

	builderConn *grpc.ClientConn
	ssgenConn   *grpc.ClientConn

	config Config
}

func New(c Config) (*Service, error) {
	s := &Service{
		server: grpc.NewServer(),
		config: c,
	}

	api.RegisterPingServer(s.server, &PingService{})

	grpcweb.WrapServer(s.server)

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

	return eg.Wait()
}

func (s *Service) builder() builderApi.BuilderServiceClient {
	return builderApi.NewBuilderServiceClient(s.builderConn)
}

func (s *Service) ssgen() ssgenApi.StaticSiteGenServiceClient {
	return ssgenApi.NewStaticSiteGenServiceClient(s.ssgenConn)
}
