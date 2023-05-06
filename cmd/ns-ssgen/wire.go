//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewControllerSSGenServiceClient,
		usecase.NewStaticSiteServerService,
		repository.New,
		repository.NewApplicationRepository,
		repository.NewBuildRepository,
		staticserver.NewBuiltIn,
		provideWebServerPort,
		provideWebServerDocumentRootPath,
		provideStorage,
		wire.FieldsOf(new(Config), "Controller", "DB", "Storage"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
