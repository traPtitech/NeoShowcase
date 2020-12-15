package httpserver

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
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

// paramAppMiddleware リクエストURLの:appIdパラメーターからAppをロードするミドルウェア
func paramAppMiddleware(am appmanager.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			appID := getRequestParamAppId(c)
			app, err := am.GetApp(appID)
			if err != nil {
				if err == appmanager.ErrNotFound {
					return echo.NewHTTPError(http.StatusNotFound)
				}

				log.WithError(err).Errorf("internal error")
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			c.Set(contextParamApp, app)
			return next(c)
		}
	}
}
