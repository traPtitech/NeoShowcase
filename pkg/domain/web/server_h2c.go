package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type H2CConfig struct {
	Port       int
	SetupRoute func(mux *http.ServeMux)
}

type H2CServer struct {
	port     int
	mux      *http.ServeMux
	listener net.Listener
}

func NewH2CServer(c H2CConfig) *H2CServer {
	mux := http.NewServeMux()
	c.SetupRoute(mux)
	return &H2CServer{port: c.Port, mux: mux}
}

func (s *H2CServer) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	s.listener = listener
	return http.Serve(listener, s.mux)
}

func (s *H2CServer) Shutdown(_ context.Context) error {
	return s.listener.Close()
}
