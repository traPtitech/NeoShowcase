package memberchecker

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const pubkeyPEM = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAraewUw7V1hiuSgUvkly9
X+tcIh0e/KKqeFnAo8WR3ez2tA0fGwM+P8sYKHIDQFX7ER0c+ecTiKpo/Zt/a6AO
gB/zHb8L4TWMr2G4q79S1gNw465/SEaGKR8hRkdnxJ6LXdDEhgrH2ZwIPzE0EVO1
eFrDms1jS3/QEyZCJ72oYbAErI85qJDF/y/iRgl04XBK6GLIW11gpf8KRRAh4vuh
g5/YhsWUdcX+uDVthEEEGOikSacKZMFGZNi8X8YVnRyWLf24QTJnTHEv+0EStNrH
HnxCPX0m79p7tBfFC2ha2OYfOtA+94ZfpZXUi2r6gJZ+dq9FWYyA0DkiYPUq9QMb
OQIDAQAB
-----END PUBLIC KEY-----
`

var pubkey *rsa.PublicKey

func (s *Service) authenticate(c echo.Context) error {
	panic("TODO")
}

func ExtractTrapExtTokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("traP_ext_token")
	if err != nil {
		return ""
	}
	return c.Value
}

func AuthorizeMemberByToken(w http.ResponseWriter, r *http.Request) {
	tokenString := ExtractTrapExtTokenFromCookie(r)
	if len(tokenString) == 0 {
		Unauthorized(w, r)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		pubkey, _ = jwt.ParseRSAPublicKeyFromPEM([]byte(pubkeyPEM))
		return pubkey, nil
	})

	if err != nil || !token.Valid {
		Unauthorized(w, r)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if nameI, ok := claims["name"]; ok {
		w.Header().Set("X-Showcase-User", nameI.(string))
	} else {
		Unauthorized(w, r)
	}
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("type")
	switch strings.ToLower(q) {
	case "soft":
		w.Header().Set("X-Showcase-User", "-")
		w.WriteHeader(200)
	case "hard":
		w.WriteHeader(403)
	default:
		w.WriteHeader(403)
	}
}
