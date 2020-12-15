package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
)

const (
	contextRequestUserID = "__req_user_id"
	contextParamApp      = "__req_param_app"
)

// getRequestUserID リクエストユーザーのIDを取得
func getRequestUserID(c echo.Context) string {
	return c.Get(contextRequestUserID).(string)
}

// getRequestParamAppId リクエストパスの:appIdパラメーターを取得
func getRequestParamAppId(c echo.Context) string {
	return c.Param("appId")
}

// getRequestParamApp リクエストパスの:appIdのappmanager.Appを取得
func getRequestParamApp(c echo.Context) appmanager.App {
	return c.Get(contextParamApp).(appmanager.App)
}
