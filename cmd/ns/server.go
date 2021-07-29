package main

import (
	"context"
	"database/sql"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/broker"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	webserver           *web.Server
	db                  *sql.DB
	builderConn         *grpc.BuilderServiceClientConn
	ssgenConn           *grpc.StaticSiteServiceClientConn
	backend             domain.Backend
	bus                 domain.Bus
	builderEventsBroker broker.BuilderEventsBroker
	cdService           usecase.ContinuousDeploymentService
}

func (s *Server) Start(ctx context.Context) error {
	go s.builderEventsBroker.Run()
	go s.cdService.Run()
	return s.webserver.Start(ctx)
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

	return eg.Wait()
}
