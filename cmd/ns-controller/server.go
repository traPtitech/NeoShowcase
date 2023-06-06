package main

import (
	"context"
	"database/sql"

	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cdservice"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
)

type Server struct {
	controllerServer *controllerServer

	db             *sql.DB
	backend        domain.Backend
	sshServer      sshserver.SSHServer
	webhook        *webhook.Receiver
	cdService      cdservice.Service
	fetcherService repofetcher.Service
	cleanerService cleaner.Service
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.backend.Start(ctx)
	})
	eg.Go(func() error {
		return s.sshServer.Start()
	})
	eg.Go(func() error {
		return s.webhook.Start(ctx)
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
		return s.controllerServer.Start(ctx)
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
		return s.sshServer.Close()
	})
	eg.Go(func() error {
		return s.webhook.Shutdown(ctx)
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
		return s.controllerServer.Shutdown(ctx)
	})

	return eg.Wait()
}
