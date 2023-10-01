package authdev

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e    *echo.Echo
	port int
}

func NewServer(header string, port int, user string) *Server {
	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set(header, user)
		return c.NoContent(http.StatusOK)
	})
	return &Server{e: e, port: port}
}

func (s *Server) Start(_ context.Context) error {
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
