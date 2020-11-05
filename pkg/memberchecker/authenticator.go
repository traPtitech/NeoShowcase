package memberchecker

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func ExtractTokenFromCookie(c echo.Context) string {
	token, err := c.Cookie("traP_ext_token")
	if err != nil {
		return ""
	}
	return token.Value
}

func (s *Service) authenticate(c echo.Context) error {
	tokenString := ExtractTokenFromCookie(c)
	if len(tokenString) == 0 {
		return unauthorized(c)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return s.pubkey, nil
	})

	if err != nil || !token.Valid {
		return unauthorized(c)

	}

	claims := token.Claims.(jwt.MapClaims)

	if nameI, ok := claims["name"]; ok {
		c.Response().Header().Set("X-Showcase-User", nameI.(string))
	} else {
		return unauthorized(c)
	}
	return c.String(http.StatusOK, "")
}

func unauthorized(c echo.Context) error {
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
