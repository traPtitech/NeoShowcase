package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetApp GET /apps/:appId
func (s *Server) GetApp(c echo.Context) error {
	app := getRequestParamApp(c)
	return c.JSON(http.StatusOK, &AppDetail{Id: app.GetID()}) // TODO
}
