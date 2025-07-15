package controller

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/webhook"
	"github.com/traPtitech/neoshowcase/pkg/usecase/cleaner"
	commitfetcher "github.com/traPtitech/neoshowcase/pkg/usecase/commit-fetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/repofetcher"
	"github.com/traPtitech/neoshowcase/pkg/usecase/sshserver"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
)

type APIServer struct {
	*web.H2CServer
}

type Server struct {
	APIServer *APIServer

	DB             *sql.DB
	Cluster        *discovery.Cluster
	Backend        domain.Backend
	SSHServer      sshserver.SSHServer
	Webhook        *webhook.Receiver
	CDService      domain.CDService
	CommitFetcher  commitfetcher.Service
	FetcherService repofetcher.Service
	CleanerService cleaner.Service
}

func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Cluster.Start(ctx)
	})
	eg.Go(func() error {
		return s.Backend.Start(ctx)
	})
	eg.Go(func() error {
		return s.SSHServer.Start()
	})
	eg.Go(func() error {
		return s.Webhook.Start(ctx)
	})
	eg.Go(func() error {
		s.CDService.Run()
		return nil
	})
	eg.Go(func() error {
		s.CommitFetcher.Run()
		return nil
	})
	eg.Go(func() error {
		s.FetcherService.Run()
		return nil
	})
	eg.Go(func() error {
		return s.CleanerService.Start(ctx)
	})
	eg.Go(func() error {
		return s.APIServer.Start(ctx)
	})

	return eg.Wait()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.DB.Close()
	})
	eg.Go(func() error {
		return s.Backend.Dispose(ctx)
	})
	eg.Go(func() error {
		return s.SSHServer.Close()
	})
	eg.Go(func() error {
		return s.Webhook.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.CDService.Stop(ctx)
	})
	eg.Go(func() error {
		return s.CommitFetcher.Stop(ctx)
	})
	eg.Go(func() error {
		return s.FetcherService.Stop(ctx)
	})
	eg.Go(func() error {
		return s.CleanerService.Shutdown(ctx)
	})
	eg.Go(func() error {
		return s.APIServer.Shutdown(ctx)
	})

	return eg.Wait()
}
