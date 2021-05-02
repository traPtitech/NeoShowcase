package webhook

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	event2 "github.com/traPtitech/neoshowcase/pkg/domain/event"
)

func (r *Receiver) githubHandler(c echo.Context) error {
	if c.Request().Header.Get("X-GitHub-Event") != "push" {
		return c.NoContent(http.StatusBadRequest)
	}
	// b: []byte形式のbody
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	body := &PushEvent{}
	if err := json.Unmarshal(b, body); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	repoURL := body.Repo.CloneURL
	secrets, err := r.keyLoader.GetWebhookSecretKeys(repoURL)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	signature := c.Request().Header.Get("X-Hub-Signature-256")
	parts := strings.Split(signature, "=")
	if len(parts) < 2 {
		return c.NoContent(http.StatusBadRequest)
	}
	// シグネチャの検証(secretsのうちひとつでもtrueならOK)
	vS := false
	for _, secret := range secrets {
		vS = vS || verifySignature(sha256.New, b, []byte(secret), parts[1])
	}
	if !vS {
		return c.NoContent(http.StatusBadRequest)
	}
	branch := strings.TrimPrefix(body.Ref, "refs/")
	r.bus.Publish(hub.Message{
		Name: event2.WebhookRepositoryPush,
		Fields: hub.Fields{
			"repository_url": repoURL,
			"branch":         branch,
		},
	})
	return c.NoContent(http.StatusNoContent)
}
