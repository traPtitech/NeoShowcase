package gateway

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
)

type APIServer struct {
	*web.H2CServer
}

type Server struct {
	APIServer *APIServer
	DB        *sql.DB
}

// Start migrates applications left without a group and then serves the API.
func (s *Server) Start(ctx context.Context) error {
	// The previous version creates applications without a group until it stops, possibly after ns-migrate has run.
	if err := repository.MigrateApplicationGroups(ctx, s.DB); err != nil {
		return err
	}
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
