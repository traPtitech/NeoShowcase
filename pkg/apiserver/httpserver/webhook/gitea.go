package webhook

import (
	"github.com/labstack/echo/v4"
)

func (r *Receiver) giteaHandler(c echo.Context) error {
	panic("implemented me") // TODO

	// 正規のPUSHイベントだったら、以下の内部イベントを発生させて、204 NoContentを返す
	// r.bus.Publish(hub.Message{
	// 	Name: event.WebhookRepositoryPush,
	// 	Fields: hub.Fields{
	// 		"repository_url": "リポジトリのURL",
	// 		"branch":         "プッシュがあったブランチ名",
	// 	},
	// })
}
