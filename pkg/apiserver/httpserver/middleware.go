package httpserver

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// debugMiddleware デバッグ用ミドルウェア
func debugMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextRequestUserID, "__DEBUG")
			return next(c)
		}
	}
}

// authenticateMiddleware TODO ユーザー認証ミドルウェア
func authenticateMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusNotImplemented, "authenticator is not implemented")
		}
	}
}
