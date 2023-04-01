package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type Server struct {
	appServer       *webAppServer
	componentServer *webComponentServer

	db             *sql.DB
	backend        domain.Backend
	bus            domain.Bus
	cdService      usecase.ContinuousDeploymentService
	fetcherService usecase.RepositoryFetcherService
	cleanerService usecase.CleanerService
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.backend.Start(ctx)
	})
	eg.Go(func() error {
		s.cdService.Run()
		return nil
	})
	eg.Go(func() error {
		s.fetcherService.Run()
		return nil
	})
	eg.Go(func() error {
		return s.cleanerService.Start(ctx)
	})
	eg.Go(func() error {
		return s.appServer.Start(ctx)
	})
	eg.Go(func() error {
		return s.componentServer.Start(ctx)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.backend.Dispose(ctx)
	})
	eg.Go(func() error {
		return s.bus.Close(ctx)
	})
	eg.Go(func() error {
		return s.cdService.Stop(ctx)
	})
	eg.Go(func() error {
		return s.fetcherService.Stop(ctx)
	})
	eg.Go(func() error {
		return s.cleanerService.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.appServer.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.componentServer.Shutdown(ctx)
	})

	return eg.Wait()
}
