package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/coreos/go-oidc/v3/oidc"
	oidc2 "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/oidc"
)

type LoginHandler web.Handler

type loginHandler struct {
	clientID     string
	clientSecret string
}

// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのログイン用ハンドラを返す関数
func NewLoginHandler(clientID, clientSecret string) LoginHandler {
	return &loginHandler{clientID: clientID, clientSecret: clientSecret}
}

func (l *loginHandler) HandleRequest(c web.Context) {
	config, _, err := oidc.NewGoogleOIDCProvider(r.Context(), l.clientID, l.clientSecret)
	if err != nil {
		return
	}

	state := randomString(64)
	nonce := randomString(64)

	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)

	http.Redirect(w, r, config.AuthCodeURL(state, oidc2.Nonce(nonce)), http.StatusFound)
}

type CallbackHandler web.Handler

type callbackHandler struct {
	clientID     string
	clientSecret string
}

// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのコールバック用ハンドラを返す関数
func newCallbackHandler(clientID, clientSecret string) CallbackHandler {
	return &callbackHandler{clientID: clientID, clientSecret: clientSecret}
}

func (c *callbackHandler) HandleRequest(c web.Context) {
	config, verifier, err := oidc.NewGoogleOIDCProvider(r.Context(), c.clientID, c.clientSecret)
	if err != nil {
		return
	}

	state, err := r.Cookie("state")
	if err != nil {
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		return
	}

	token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return
	}
	idToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		return
	}

	if idToken.Nonce != nonce.Value {
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage
	}{token, new(json.RawMessage)}

	if err = idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}

	http.SetCookie(w, c)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
