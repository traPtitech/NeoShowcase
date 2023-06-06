//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewControllerSSGenServiceClient,
		ssgen.NewGeneratorService,
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
