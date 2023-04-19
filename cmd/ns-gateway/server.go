package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	appServer *gatewayServer
	db        *sql.DB
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

	return eg.Wait()
}
