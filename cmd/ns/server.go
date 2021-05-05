package main

import (
	"context"
	"database/sql"

	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/broker"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	webserver           *web.Server
	db                  *sql.DB
	builderConn         *grpc.BuilderServiceClientConn
	ssgenConn           *grpc.StaticSiteServiceClientConn
	backend             backend.Backend
	appmanager          appmanager.Manager
	bus                 eventbus.Bus
	builderEventsBroker broker.BuilderEventsBroker
}

func (s *Server) Start(ctx context.Context) error {
	go s.builderEventsBroker.Run()
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
		return s.appmanager.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.bus.Close(ctx)
	})

	return eg.Wait()
}
