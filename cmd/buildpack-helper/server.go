package buildpackhelper

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type APIServer struct {
	*web.H2CServer
}

type Server struct {
	Helper *APIServer
}

func (s *Server) Start(ctx context.Context) error {
	return s.Helper.Start(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Helper.Shutdown(ctx)
}
