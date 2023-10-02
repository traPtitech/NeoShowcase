package ssgen

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
)

type Server struct {
	DB      *sql.DB
	Service ssgen.GeneratorService
	Health  healthcheck.Server
	Engine  domain.StaticServer
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Service.Start(ctx)
	})
	eg.Go(func() error {
		return s.Health.Start(ctx)
	})
	eg.Go(func() error {
		return s.Engine.Start(ctx)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.DB.Close()
	})
	eg.Go(func() error {
		return s.Service.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.Health.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.Engine.Shutdown(ctx)
	})

	return eg.Wait()
}
