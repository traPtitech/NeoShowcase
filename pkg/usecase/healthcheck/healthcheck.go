package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
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
	server *http.Server
}

func NewServer(port Port, fn Func) Server {
	e := echo.New()
	e.GET("/healthz", func(c *echo.Context) error {
		if fn() {
			return c.NoContent(http.StatusOK)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	})

	return &server{
		server: &http.Server{
			Addr:        fmt.Sprintf(":%d", int(port)),
			Handler:     e,
			ReadTimeout: 30 * time.Second,
		},
	}
}

func (h *server) Start(_ context.Context) error {
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (h *server) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
