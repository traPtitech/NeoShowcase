package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/event"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type WebhookReceiverHandler web.Handler

type webhookReceiverHandler struct {
	eventbus domain.Bus
	verifier usecase.GitPushWebhookService
}

func NewWebhookReceiverHandler(eventbus domain.Bus, verifier usecase.GitPushWebhookService) WebhookReceiverHandler {
	return &webhookReceiverHandler{
		eventbus: eventbus,
		verifier: verifier,
	}
}

func (h *webhookReceiverHandler) HandleRequest(c web.Context) error {
	var (
		repoURL string
		branch  string
		err     error
	)
	switch {
	case c.Request().Header.Get("X-Gitea-Delivery") != "":
		repoURL, branch, err = h.extractFromGitea(c)
	case c.Request().Header.Get("X-GitHub-Delivery") != "":
		repoURL, branch, err = h.extractFromGitHub(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err != nil {
		return err
	}

	h.eventbus.Publish(event.WebhookRepositoryPush, domain.Fields{
		"repository_url": repoURL,
		"branch":         branch,
	})
	return c.NoContent(http.StatusNoContent)
}

func (h *webhookReceiverHandler) extractFromGitea(c web.Context) (string, string, error) {
	if c.Request().Header.Get("X-Gitea-Event") != "push" {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	var body struct {
		Ref  string `json:"ref"`
		Repo struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
			CloneURL string `json:"clone_url"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(b, &body); err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	repoURL := body.Repo.CloneURL
	signature := c.Request().Header.Get("X-Gitea-Signature")
	branch := strings.TrimPrefix(body.Ref, "refs/")

	valid, err := h.verifier.VerifySignature(c.Request().Context(), repoURL, signature, b)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !valid {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}
	name := body.Repo.Name
	owner := body.Repo.Owner.Login

	exists, err := h.verifier.CheckRepositoryExists(c.Request().Context(), repoURL, owner, name)

	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !exists {
		return "", "", echo.NewHTTPError(http.StatusNotFound)
	}
	return repoURL, branch, nil
}

func (h *webhookReceiverHandler) extractFromGitHub(c web.Context) (string, string, error) {
	if c.Request().Header.Get("X-GitHub-Event") != "push" {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	var body struct {
		Ref  string `json:"ref"`
		Repo struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
			CloneURL string `json:"clone_url"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(b, &body); err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	repoURL := body.Repo.CloneURL
	parts := strings.Split(c.Request().Header.Get("X-Hub-Signature-256"), "=")
	if len(parts) < 2 {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}
	branch := strings.TrimPrefix(body.Ref, "refs/")

	valid, err := h.verifier.VerifySignature(c.Request().Context(), repoURL, parts[1], b)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !valid {
		return "", "", echo.NewHTTPError(http.StatusBadRequest)
	}

	name := body.Repo.Name
	owner := body.Repo.Owner.Login

	exists, err := h.verifier.CheckRepositoryExists(c.Request().Context(), repoURL, owner, name)

	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !exists {
		return "", "", echo.NewHTTPError(http.StatusNotFound)
	}
	return repoURL, branch, nil
}
