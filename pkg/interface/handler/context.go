package handler

import "github.com/labstack/echo/v4"

type Context interface {
	echo.Context
	CookieValue(name string) string
}
