package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetDebug1(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) GetDebug2(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
