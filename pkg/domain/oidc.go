package domain

import (
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/oidc"
	"math/rand"
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

	l := NewLogin()
	config, verifier, err := oidc.NewGoogleOIDCProvider(r.Context(), l.clientID, l.clientSecret)
	if err != nil {
		return nil
	}

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

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
