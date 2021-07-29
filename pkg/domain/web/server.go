package web

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	conf Config
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
	conf.Router.SetupRoute(e)

	return &Server{echo: e, conf: conf}
}

func (s *Server) Start(_ context.Context) error {
	return s.echo.Start(fmt.Sprintf(":%d", s.conf.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
