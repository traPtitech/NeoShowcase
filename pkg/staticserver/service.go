package staticserver

import (
	"github.com/labstack/echo/v4"
)

type Service struct {
	root *echo.Echo
}

func (s *Service) Run() error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Any("/*", s.serveSites)

	s.root = e
	return e.Start(":8888")
}

func (s *Service) serveSites(c echo.Context) error {
	return nil
}
