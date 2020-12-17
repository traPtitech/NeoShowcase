package webhook

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/event"
)

func (r *Receiver) giteaHandler(c echo.Context) error {
	if c.Request().Header.Get("X-Gitea-Event") != "push" {
		return errors.New("Webhook event isn't push event")
	}
	var body PushEvent
	if err := c.Bind(body); err != nil {
		return errors.New("Couldn't bind request body")
	}
	repoURL := body.Repo.HTMLURL
	secrets, err := r.keyLoader.GetWebhookSecretKeys(repoURL)
	if err != nil {
		return errors.New("Couldn't get webhook secret keys")
	}
	signature := c.Request().Header.Get("X-Gitea-Signature")
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.New("Couldn't read request body")
	}
	if !verifySignature([]byte(signature), secrets[0], b) {
		return errors.New("Invalid signature")
	}
	branch := strings.Trim(body.Ref, "refs/heads/")
	r.bus.Publish(hub.Message{
		Name: event.WebhookRepositoryPush,
		Fields: hub.Fields{
			"repository_url": repoURL,
			"branch":         branch,
		},
	})
	return nil
}
