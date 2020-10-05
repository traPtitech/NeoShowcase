package memberchecker

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Service struct {
	echo   *echo.Echo
	pubkey *rsa.PublicKey
	config Config
}

func New(c Config) (*Service, error) {
	s := &Service{
		config: c,
	}

	// router初期化
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Debug = false
	e.Use(middleware.Recover())
	e.GET("/", s.authenticate)
	s.echo = e

	// JWT公開鍵をパース
	pubkey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(c.getJWTPublicKeyPEM()))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT RSA public key from pem: %w", err)
	}
	s.pubkey = pubkey

	return s, nil
}

func (s *Service) Start(_ context.Context) error {
	return s.echo.Start(fmt.Sprintf(":%d", s.config.getHTTPPort()))
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
