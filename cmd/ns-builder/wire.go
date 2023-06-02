//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewControllerBuilderServiceClient,
		usecase.NewBuilderService,
		repository.New,
		repository.NewApplicationRepository,
		repository.NewArtifactRepository,
		repository.NewBuildRepository,
		repository.NewGitRepositoryRepository,
		provideBuildpackBackend,
		provideStorage,
		initBuildkitClient,
		provideRepositoryPublicKey,
		wire.FieldsOf(new(Config), "Controller", "DB", "Storage", "Image"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
