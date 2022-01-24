package handler

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"net/http"
)

// TODO: signup & login
type AuthHandler Handler

type authHandler struct {
}

func (a *authHandler) HandleRequest(c Context) error {
	lh := domain.NewLoginHandler()

	ch := domain.NewCallbackHandler()

	return c.NoContent(http.StatusOK)
}
