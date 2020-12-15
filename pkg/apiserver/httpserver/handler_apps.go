package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetApp GET /apps/:appId
func (s *Server) GetApp(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, &AppDetail{Id: getRequestParamAppId(c)})
}
