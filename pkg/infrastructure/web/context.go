package web

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
)

type Context struct {
	echo.Context
}

func (c *Context) CookieValue(name string) string {
	token, err := c.Cookie(name)
	if err != nil {
		return ""
	}
	return token.Value
}

func UnwrapHandler(h handler.Handler) echo.HandlerFunc {
	return UnwrapHandlerFunc(h.HandleRequest)
}

func UnwrapHandlerFunc(f func(ctx handler.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return f(c.(*Context))
	}
}

func wrapContext(c echo.Context) *Context {
	return &Context{Context: c}
}

func wrapContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(wrapContext(c))
		}
	}
}
