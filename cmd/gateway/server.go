package gateway

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type APIServer struct {
	*web.H2CServer
}

type Server struct {
	APIServer *APIServer
	DB        *sql.DB
}

func (s *Server) Start(ctx context.Context) error {
	return s.APIServer.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.APIServer.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.DB.Close()
	})

	return eg.Wait()
}
