package webhook

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
)

type PushEvent struct {
	Secret string              `json:"secret"`
	Ref    string              `json:"ref"`
	Repo   PushEventRepository `json:"repository"`
}

type PushEventRepository struct {
	HTMLURL string `json:"html_url"`
}

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

// Validate checks the hmac signature of the message
// using a hex encoded signature.
func verifySignature(h func() hash.Hash, message, key []byte, signature string) bool {
	decoded, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return validate(h, message, key, decoded)
}

func validate(h func() hash.Hash, message, key, signature []byte) bool {
	mac := hmac.New(h, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}
