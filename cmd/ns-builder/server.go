package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	buildkit "github.com/moby/buildkit/client"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	igrpc "github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	db         *sql.DB
	eventbus   domain.Bus
	grpcServer *grpc.Server
	buildkit   *buildkit.Client

	port    igrpc.TCPListenPort
	builder *igrpc.BuilderService
}

func (s *Server) Start(ctx context.Context) error {
	pb.RegisterBuilderServiceServer(s.grpcServer, s.builder)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return s.grpcServer.Serve(listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		s.grpcServer.GracefulStop()
		return nil
	})
	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.buildkit.Close()
	})
	eg.Go(func() error {
		return s.eventbus.Close(ctx)
	})

	return eg.Wait()
}
