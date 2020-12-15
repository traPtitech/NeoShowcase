package webhook

import (
	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	"net/http"
)

// Receiver Webhookレシーバー
type Receiver struct {
	bus       *hub.Hub
	keyLoader SecretKeyLoader
}

// SecretKeyLoader Webhookシグネチャのシークレットキーを取得する
type SecretKeyLoader interface {
	// GetWebhookSecretKeys 指定したリポジトリのWebhookシークレットキーを取得します
	GetWebhookSecretKeys(repositoryUrl string) ([]string, error)
}

func NewReceiver(bus *hub.Hub, keyLoader SecretKeyLoader) *Receiver {
	return &Receiver{bus: bus, keyLoader: keyLoader}
}

// Handler Gitホスティングサービスから送られてきたWebhookを受け取るハンドラ
func (r *Receiver) Handler(c echo.Context) error {
	switch {
	case c.Request().Header.Get("X-Gitea-Delivery") != "":
		return r.giteaHandler(c) // GiteaからのWebhook
	case c.Request().Header.Get("X-GitHub-Delivery") != "":
		return r.githubHandler(c) // GitHubからのWebhook
	default:
		return c.NoContent(http.StatusBadRequest)
	}
}
