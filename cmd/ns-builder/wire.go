//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewComponentServiceClient,
		usecase.NewBuilderService,
		repository.NewArtifactRepository,
		repository.NewBuildRepository,
		repository.NewGitRepositoryRepository,
		admindb.New,
		provideStorage,
		initBuildkitClient,
		provideRepositoryPublicKey,
		wire.FieldsOf(new(Config), "NS", "DB", "Storage"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
