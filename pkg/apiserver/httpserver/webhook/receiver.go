package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
)

type PushEvent struct {
	Ref  string              `json:"ref,omitempty"`
	Repo PushEventRepository `json:"repository,omitempty"`
}

type PushEventRepository struct {
	HTMLURL string `json:"html_url,omitempty"`
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
func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}
