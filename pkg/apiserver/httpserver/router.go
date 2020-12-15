package httpserver

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e      *echo.Echo
	config Config
}

type Config struct {
	Debug bool
	Port  int
}

func New(config Config) *Server {
	s := &Server{
		e:      echo.New(),
		config: config,
	}
	e := s.e
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	api := e.Group("")
	if config.Debug {
		api.Use(debugMiddleware())
	} else {
		api.Use(authenticateMiddleware())
	}

	api.GET("/apps/:appId", s.GetApp)

	return s
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.config.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
