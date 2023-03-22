package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/broker"
	igrpc "github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type Server struct {
	grpcServer *grpc.Server
	grpcPort   igrpc.TCPListenPort
	appService *igrpc.ApplicationService

	webserver           *web.Server
	db                  *sql.DB
	builderConn         *igrpc.BuilderServiceClientConn
	ssgenConn           *igrpc.StaticSiteServiceClientConn
	backend             domain.Backend
	bus                 domain.Bus
	builderEventsBroker broker.BuilderEventsBroker
	cdService           usecase.ContinuousDeploymentService
	fetcherService      usecase.RepositoryFetcherService
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.backend.Start(ctx)
	})
	eg.Go(func() error {
		return s.builderEventsBroker.Run()
	})
	eg.Go(func() error {
		s.cdService.Run()
		return nil
	})
	eg.Go(func() error {
		s.fetcherService.Run()
		return nil
	})
	eg.Go(func() error {
		return s.webserver.Start(ctx)
	})
	eg.Go(func() error {
		pb.RegisterApplicationServiceServer(s.grpcServer, s.appService)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.grpcPort))
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return s.grpcServer.Serve(listener)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.webserver.Shutdown(ctx)
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
		return s.backend.Dispose(ctx)
	})
	eg.Go(func() error {
		return s.bus.Close(ctx)
	})
	eg.Go(func() error {
		return s.cdService.Stop(ctx)
	})
	eg.Go(func() error {
		return s.fetcherService.Stop(ctx)
	})

	return eg.Wait()
}
