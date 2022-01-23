package domain

import (
	"net/http"
	"os"
	"time"
)

type LoginHandler struct {
	clientID     string
	clientSecret string
}

func NewLogin() *LoginHandler {
	clientID := os.Getenv("clientID")
	clientSecret := os.Getenv("clientSecret")
	return &LoginHandler{clientID: clientID, clientSecret: clientSecret}
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewLoginHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのログイン用ハンドラを返す関数
	return http.NotFoundHandler()
}

func NewCallbackHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのコールバック用ハンドラを返す関数
	return http.NotFoundHandler()
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
