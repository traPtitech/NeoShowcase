package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
)

type Server struct {
	db     *sql.DB
	svc    ssgen.GeneratorService
	engine domain.StaticServer
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.svc.Start(ctx)
	})
	eg.Go(func() error {
		return s.engine.Start(ctx)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.svc.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.engine.Shutdown(ctx)
	})

	return eg.Wait()
}
