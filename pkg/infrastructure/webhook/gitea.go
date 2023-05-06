package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type giteaPushPayload struct {
	Repository struct {
		HTMLURL  string `json:"html_url"`
		SSHURL   string `json:"ssh_url"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
}

func (r *Receiver) giteaHandler(c echo.Context) error {
	var p giteaPushPayload
	err := json.NewDecoder(c.Request().Body).Decode(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// https://docs.gitea.io/en-us/usage/webhooks/
	urls := []string{
		p.Repository.HTMLURL,  // http://localhost:3000/gitea/webhooks
		p.Repository.SSHURL,   // ssh://gitea@localhost:2222/gitea/webhooks.git
		p.Repository.CloneURL, // http://localhost:3000/gitea/webhooks.git
	}
	go r.updateURLs(urls)

	return c.NoContent(http.StatusOK)
}
