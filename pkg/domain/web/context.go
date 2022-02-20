package web

import (
	"github.com/labstack/echo/v4"
)

type Context interface {
	echo.Context
	CookieValue(name string) (string, error)
}

type DefaultContext struct {
	echo.Context
}

func (c *DefaultContext) CookieValue(name string) (string, error) {
	token, err := c.Cookie(name)
	return token.Value, err
}

func UnwrapHandler(h Handler) echo.HandlerFunc {
	return UnwrapHandlerFunc(h.HandleRequest)
}

func UnwrapHandlerFunc(f func(ctx Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return f(c.(*DefaultContext))
	}
}

func wrapContext(c echo.Context) Context {
	return &DefaultContext{Context: c}
}

func WrapContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(wrapContext(c))
		}
	}
}
