package webhook

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/go-playground/webhooks/v6/gitea"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

var giteaHook = lo.Must(gitea.New())

func (r *Receiver) giteaHandler(c echo.Context) error {
	rawPayload, err := giteaHook.Parse(c.Request(), gitea.PushEvent, gitea.RepositoryEvent)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// https://docs.gitea.io/en-us/usage/webhooks/
	switch p := rawPayload.(type) {
	case gitea.PushPayload:
		urls := []string{
			p.Repo.HTMLURL,  // http://localhost:3000/gitea/webhooks
			p.Repo.SSHURL,   // ssh://gitea@localhost:2222/gitea/webhooks.git
			p.Repo.CloneURL, // http://localhost:3000/gitea/webhooks.git
		}
		go r.updateURLs(urls)
	case gitea.RepositoryPayload:
		slog.Info("Repository event received", "action", p.Action)
		if err := r.giteaIntegration.Sync(context.Background()); err != nil {
			slog.Error("Failed to sync gitea repositories", "error", err)
		}
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("unsupported payload type"))
	}

	return c.NoContent(http.StatusOK)
}
