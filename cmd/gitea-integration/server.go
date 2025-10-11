package giteaintegration

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

type APIServer struct {
	*web.H2CServer
}

type Server struct {
	Integration *giteaintegration.Integration
	DB          *sql.DB
	APIServer   *APIServer
}

func (s *Server) Start(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(s.Integration.Start)
	eg.Go(func() error { return s.APIServer.Start(ctx) })
	return eg.Wait()
}

func (s *Server) Shutdown(_ context.Context) error {
	var eg errgroup.Group

	eg.Go(func() error {
		return s.Integration.Shutdown()
	})
	eg.Go(func() error {
		return s.DB.Close()
	})

	return eg.Wait()
}
