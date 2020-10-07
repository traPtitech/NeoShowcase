package apiserver

import (
	"github.com/labstack/echo/v4"
)

type Service struct {
	echo *echo.Echo
}

func New() *Service {
	s := &Service{
		echo: nil,
	}

	return s
}
