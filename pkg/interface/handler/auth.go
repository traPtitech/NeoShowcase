package handler

import "net/http"

// TODO: signup & login
type AuthHandler Handler

type authHandler struct {
	
}

func (a *authHandler) HandleRequest(c Context) error {

	return c.NoContent(http.StatusOK)
}

