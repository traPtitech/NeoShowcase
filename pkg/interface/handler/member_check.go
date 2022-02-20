package handler

import (
	"net/http"
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type MemberCheckHandler web.Handler

type memberCheckHandler struct {
	s usecase.MemberCheckService
}

func NewMemberCheckHandler(s usecase.MemberCheckService) MemberCheckHandler {
	return &memberCheckHandler{s: s}
}

func (h *memberCheckHandler) HandleRequest(c web.Context) error {
	unauthorized := func() error {
		q := c.QueryParam("type")
		switch strings.ToLower(q) {
		case "soft":
			c.Response().Header().Set("X-Showcase-User", "-")
			return c.String(http.StatusOK, "")
		case "hard":
			return c.NoContent(http.StatusForbidden)
		default:
			return c.NoContent(http.StatusForbidden)
		}
	}

	tokenString, err := c.CookieValue("traP_ext_token")
	if len(tokenString) == 0 {
		return unauthorized()
	}

	id, err := h.s.Check(tokenString)
	if err != nil {
		return unauthorized()
	}

	c.Response().Header().Set("X-Showcase-User", id)
	return c.String(http.StatusOK, "")
}
