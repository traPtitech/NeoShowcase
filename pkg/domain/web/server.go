package web

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port       int
	SetupRoute func(e *echo.Echo)
}

type Server struct {
	echo *echo.Echo
	port int
}

func NewServer(conf Config) *Server {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Debug = false

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	e.Use(WrapContextMiddleware())
	conf.SetupRoute(e)

	return &Server{echo: e, port: conf.Port}
}

func (s *Server) Start(_ context.Context) error {
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
