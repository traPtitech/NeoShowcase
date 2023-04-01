package web

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type H2CConfig struct {
	Port       int
	SetupRoute func(mux *http.ServeMux)
}

type H2CServer struct {
	server *http.Server
}

func NewH2CServer(c H2CConfig) *H2CServer {
	mux := http.NewServeMux()
	c.SetupRoute(mux)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	return &H2CServer{server: server}
}

func (s *H2CServer) Start(_ context.Context) error {
	return s.server.ListenAndServe()
}

func (s *H2CServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
