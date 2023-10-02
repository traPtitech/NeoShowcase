package giteaintegration

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

type Server struct {
	Integration *giteaintegration.Integration
	DB          *sql.DB
}

func (s *Server) Start(_ context.Context) error {
	return s.Integration.Start()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Integration.Shutdown()
	})
	eg.Go(func() error {
		return s.DB.Close()
	})

	return eg.Wait()
}
