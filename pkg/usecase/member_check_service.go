package usecase

import (
	"crypto/rsa"

	"github.com/friendsofgo/errors"
	"github.com/golang-jwt/jwt/v4"
)

type TrapShowcaseJWTPublicKeyPEM string

type MemberCheckService interface {
	Check(token string) (traPID string, err error)
}

type memberCheckService struct {
	pubkey *rsa.PublicKey
}

func NewMemberCheckService(pem TrapShowcaseJWTPublicKeyPEM) (MemberCheckService, error) {
	// JWT公開鍵をパース
	pubkey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse JWT RSA public key from pem")
	}
	return &memberCheckService{pubkey: pubkey}, nil
}

func (s *memberCheckService) Check(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return s.pubkey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	nameI, ok := claims["name"]
	if !ok {
		return "", errors.New("invalid token")
	}
	traPID, ok := nameI.(string)
	if !ok {
		return "", errors.New("invalid token")
	}
	return traPID, nil
}
