package webhook

import (
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

var githubHook = lo.Must(github.New())

func (r *Receiver) githubHandler(c echo.Context) error {
	rawPayload, err := githubHook.Parse(c.Request(), github.PingEvent, github.PushEvent)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// https://docs.github.com/en/rest/repos/repos
	switch p := rawPayload.(type) {
	case github.PingPayload:
		// no-op
	case github.PushPayload:
		urls := []string{
			p.Repository.HTMLURL,  // https://github.com/octocat/Hello-World
			p.Repository.GitURL,   // git:github.com/octocat/Hello-World.git
			p.Repository.SSHURL,   // git@github.com:octocat/Hello-World.git
			p.Repository.CloneURL, // https://github.com/octocat/Hello-World.git
		}
		go r.updateURLs(urls)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("unsupported payload type"))
	}

	return c.NoContent(http.StatusOK)
}
