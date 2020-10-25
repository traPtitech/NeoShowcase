package memberchecker

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"strings"
)

func (s *Service) authenticate(c echo.Context) error {
	panic("TODO")
}

func ExtractTrapExtTokenFromCookie(c echo.Context) string {
	token, err := c.Cookie("traP_ext_token")
	if err != nil {
		return ""
	}
	return token.Value
}

func (s *Service) AuthorizeMemberByToken(c echo.Context) {
	tokenString := ExtractTrapExtTokenFromCookie(c)
	if len(tokenString) == 0 {
		Unauthorized(c)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return s.pubkey, nil
	})

	if err != nil || !token.Valid {
		Unauthorized(c)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if nameI, ok := claims["name"]; ok {
		c.Response().Header().Set("X-Showcase-User", nameI.(string))
	} else {
		Unauthorized(c)
	}
}

func Unauthorized(c echo.Context) {
	q := c.Request().URL.Query().Get("type")
	switch strings.ToLower(q) {
	case "soft":
		c.Response().Header().Set("X-Showcase-User", "-")
		c.Response().WriteHeader(200)
	case "hard":
		c.Response().WriteHeader(403)
	default:
		c.Response().WriteHeader(403)
	}
}
