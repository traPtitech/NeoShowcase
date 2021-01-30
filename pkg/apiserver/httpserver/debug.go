package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	"net/http"
)

func (s *Server) GetDebug1(c echo.Context) error {
	app, err := s.config.AppManager.GetApp("test")
	if err != nil {
		return err
	}
	if err := app.RequestBuild(context.Background()); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetDebug2(c echo.Context) error {
	app, err := s.config.AppManager.GetApp("test")
	if err != nil {
		return err
	}
	if err := app.Start(appmanager.AppStartArgs{}); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
