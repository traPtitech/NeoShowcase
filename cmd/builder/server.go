package builder

import (
	"context"

	buildkit "github.com/moby/buildkit/client"
	"golang.org/x/sync/errgroup"

	"github.com/traPtitech/neoshowcase/pkg/usecase/builder"
)

type Server struct {
	Buildkit *buildkit.Client
	Builder  builder.Service
}

func (s *Server) Start(ctx context.Context) error {
	return s.Builder.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Buildkit.Close()
	})
	eg.Go(func() error {
		return s.Builder.Shutdown(ctx)
	})

	return eg.Wait()
}
