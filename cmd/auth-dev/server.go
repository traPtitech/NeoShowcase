package authdev

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type Server struct {
	server *http.Server
}

func NewServer(header string, port int, user string) *Server {
	e := echo.New()
	e.Any("/*", func(c *echo.Context) error {
		c.Response().Header().Set(header, user)
		return c.NoContent(http.StatusOK)
	})
	return &Server{server: &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     e,
		ReadTimeout: 30 * time.Second,
	}}
}

func (s *Server) Start(_ context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
