package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

type Server struct {
	integration *giteaintegration.Integration
	db          *sql.DB
}

func (s *Server) Start(_ context.Context) error {
	return s.integration.Start()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.integration.Shutdown()
	})
	eg.Go(func() error {
		return s.db.Close()
	})

	return eg.Wait()
}
