package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
	"net/http"
)

func (s *Server) GetDebug1(c echo.Context) error {
	app, err := s.config.AppManager.CreateApp(appmanager.CreateAppArgs{
		Owner:         "test",
		Name:          "test",
		RepositoryURL: "https://github.com/yeasy/simple-web.git",
		BranchName:    "master",
		BuildType:     appmanager.BuildTypeImage,
	})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	envID := app.GetEnvs()[0].GetID()
	if err := app.RequestBuild(context.Background(), envID); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetDebug2(c echo.Context) error {
	app, err := s.config.AppManager.GetAppByRepository("https://github.com/yeasy/simple-web.git")
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := app.Start(appmanager.AppStartArgs{
		EnvironmentID: app.GetEnvs()[0].GetID(),
	}); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.NoContent(http.StatusNoContent)
}
