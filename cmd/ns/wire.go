//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		NewServer,
		web.NewServer,
		usecase.NewGitPushWebhookService,
		repository.NewWebhookSecretRepository,
		eventbus.NewLocal,
		admindb.New,
		handlerSet,
		provideAdminDBConfig,
		provideWebServerConfig,
		hub.New,
		wire.Struct(new(Router), "*"),
		wire.Bind(new(web.Router), new(*Router)),
	)
	return nil, nil
}
