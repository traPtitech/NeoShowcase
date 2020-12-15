package httpserver

import "github.com/labstack/echo/v4"

const (
	contextRequestUserID = "__req_user_id"
)

// getRequestUserID リクエストユーザーのIDを取得
func getRequestUserID(c echo.Context) string {
	return c.Get(contextRequestUserID).(string)
}

// getRequestParamAppId リクエストパスの:appIdパラメーターを取得
func getRequestParamAppId(c echo.Context) string {
	return c.Param("appId")
}
