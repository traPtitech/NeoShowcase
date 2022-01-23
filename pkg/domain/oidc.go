package domain

import (
	_ "github.com/coreos/go-oidc/v3/oidc"
	oidc2 "github.com/coreos/go-oidc/v3/oidc"
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

func newLoginHandler() *LoginHandler {
	clientID := os.Getenv("clientID")
	clientSecret := os.Getenv("clientSecret")
	return &LoginHandler{clientID: clientID, clientSecret: clientSecret}
}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func NewLoginHandler() http.Handler {
	// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのログイン用ハンドラを返す関数

	l := newLoginHandler()

	return l
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
