package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	igrpc "github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type Server struct {
	db         *sql.DB
	svc        usecase.StaticSiteServerService
	grpcServer *grpc.Server
	engine     domain.Engine
	sss        *igrpc.StaticSiteService
	port       igrpc.TCPListenPort
}

func (s *Server) Start(ctx context.Context) error {
	pb.RegisterStaticSiteServiceServer(s.grpcServer, s.sss)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.engine.Start(ctx)
	})
	eg.Go(func() error {
		return s.svc.Reload(ctx)
	})
	eg.Go(func() error {
		return s.grpcServer.Serve(listener)
	})

	return eg.Wait()
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
		return s.engine.Shutdown(ctx)
	})

	return eg.Wait()
}
