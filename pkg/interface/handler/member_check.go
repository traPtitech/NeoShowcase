package handler

import (
	"net/http"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

type TokenCookieName string

type MemberCheckHandler web.Handler

type memberCheckHandler struct {
	s          usecase.MemberCheckService
	cookieName string
}

func NewMemberCheckHandler(s usecase.MemberCheckService, cookieName TokenCookieName) MemberCheckHandler {
	return &memberCheckHandler{
		s:          s,
		cookieName: string(cookieName),
	}
}

func (h *memberCheckHandler) HandleRequest(c web.Context) error {
	unauthorized := func() error {
		authType := c.Request().Header.Get(web.HeaderNameAuthorizationType)
		switch authType {
		case "", "none":
			return c.String(http.StatusOK, "")
		case "soft":
			c.Response().Header().Set(web.HeaderNameShowcaseUser, "-")
			return c.String(http.StatusOK, "")
		case "hard":
			return c.NoContent(http.StatusForbidden)
		default:
			return c.String(http.StatusBadRequest, "bad auth type")
		}
	}

	tokenString, err := c.CookieValue(h.cookieName)
	if len(tokenString) == 0 {
		return unauthorized()
	}

	id, err := h.s.Check(tokenString)
	if err != nil {
		return unauthorized()
	}

	c.Response().Header().Set(web.HeaderNameShowcaseUser, id)
	return c.String(http.StatusOK, "")
}
