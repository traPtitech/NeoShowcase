package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type Server struct {
	appServer *gatewayServer
	db        *sql.DB
	bus       domain.Bus
}

func (s *Server) Start(ctx context.Context) error {
	return s.appServer.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.appServer.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.bus.Close(ctx)
	})

	return eg.Wait()
}
