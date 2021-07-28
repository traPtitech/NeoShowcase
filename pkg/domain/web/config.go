package web

import "github.com/labstack/echo/v4"

type Router interface {
	SetupRoute(e *echo.Echo)
}

type Config struct {
	Port   int
	Router Router
}
