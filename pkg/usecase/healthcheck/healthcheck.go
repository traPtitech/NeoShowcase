package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Port int
	Func func() bool
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type server struct {
	server *echo.Echo
	port   int
}

func NewServer(port Port, fn Func) Server {
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.GET("/healthz", func(c echo.Context) error {
		if fn() {
			return c.NoContent(http.StatusOK)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	})

	return &server{
		server: e,
		port:   int(port),
	}
}

func (h *server) Start(_ context.Context) error {
	return h.server.Start(fmt.Sprintf(":%v", h.port))
}

func (h *server) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
