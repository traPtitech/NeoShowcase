//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func NewServer() (*web.Server, error) {
	wire.Build(
		web.NewServer,
		usecase.NewMemberCheckService,
		handler.NewMemberCheckHandler,
		providePubKeyPEM,
		provideServerConfig,
		provideTokenCookieName,
		wire.Struct(new(Router), "*"),
		wire.Bind(new(web.Router), new(*Router)),
	)
	return nil, nil
}
