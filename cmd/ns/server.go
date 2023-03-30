package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/friendsofgo/errors"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	grpc2 "github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc/pb"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type Server struct {
	appServer        grpc2.ApplicationServiceGRPCServer
	appService       pb.ApplicationServiceServer
	componentServer  grpc2.ComponentServiceGRPCServer
	componentService domain.ComponentService

	db             *sql.DB
	backend        domain.Backend
	bus            domain.Bus
	cdService      usecase.ContinuousDeploymentService
	fetcherService usecase.RepositoryFetcherService
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.backend.Start(ctx)
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
		pb.RegisterApplicationServiceServer(s.appServer, s.appService)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GRPC.App.Port))
		if err != nil {
			return errors.Wrap(err, "failed to start app server")
		}
		return s.appServer.Serve(listener)
	})
	eg.Go(func() error {
		pb.RegisterComponentServiceServer(s.componentServer, s.componentService)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GRPC.Component.Port))
		if err != nil {
			return errors.Wrap(err, "failed to start component server")
		}
		return s.componentServer.Serve(listener)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.db.Close()
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
	eg.Go(func() error {
		s.appServer.Stop()
		return nil
	})
	eg.Go(func() error {
		s.componentServer.Stop()
		return nil
	})

	return eg.Wait()
}
