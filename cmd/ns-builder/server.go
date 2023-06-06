package main

import (
	"context"
	"database/sql"

	buildkit "github.com/moby/buildkit/client"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/usecase/builder"
)

type Server struct {
	db       *sql.DB
	buildkit *buildkit.Client
	builder  builder.Service
}

func (s *Server) Start(ctx context.Context) error {
	return s.builder.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.db.Close()
	})
	eg.Go(func() error {
		return s.buildkit.Close()
	})
	eg.Go(func() error {
		return s.builder.Shutdown(ctx)
	})

	return eg.Wait()
}
