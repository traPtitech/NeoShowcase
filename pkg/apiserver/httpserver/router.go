package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/apiserver/httpserver/webhook"
	"github.com/traPtitech/neoshowcase/pkg/appmanager"
)

type Server struct {
	e      *echo.Echo
	config Config
}

type Config struct {
	Debug      bool
	Port       int
	Bus        *hub.Hub
	AppManager appmanager.Manager
}

func New(config Config) *Server {
	s := &Server{
		e:      echo.New(),
		config: config,
	}
	e := s.e
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	api := e.Group("")
	if config.Debug {
		api.Use(debugMiddleware())
		api.GET("/_debug/1", s.GetDebug1)
		api.GET("/_debug/2", s.GetDebug2)
	} else {
		api.Use(authenticateMiddleware())
	}

	{
		apiApps := api.Group("/apps")
		{
			apiAppsAppId := apiApps.Group("/:appId", paramAppMiddleware(config.AppManager)) // TODO アクセス権限チェック
			{
				apiAppsAppId.GET("", s.GetApp)
			}
		}
	}

	apiNoAuth := e.Group("")
	apiNoAuth.POST("/_webhook", webhook.NewReceiver(config.Bus, s).Handler)

	return s
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.config.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) GetWebhookSecretKeys(repositoryUrl string) ([]string, error) {
	return []string{"__test"}, nil // TODO
	// repo, err := s.config.AppManager.GetRepoByURL(repositoryUrl)
}
