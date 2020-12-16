package webhook

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/event"
)

func (r *Receiver) githubHandler(c echo.Context) error {
	if c.Request().Header.Get("X-GitHub-Event") != "push" {
		return errors.New("Webhook event isn't push event")
	}
	var body PushEvent
	if err := c.Bind(body); err != nil {
		return errors.New("Couldn't read request body")
	}
	repoURL := body.Repo.HTMLURL
	secret, err := r.keyLoader.GetWebhookSecretKeys(repoURL)
	if err != nil {
		return errors.New("Couldn't get webhook secret keys")
	}
	signature := c.Request().Header.Get("X-GitHub-Signature")
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if !verifySignature([]byte(signature), secret[0], b) {
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
