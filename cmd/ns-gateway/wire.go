//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase/apiserver"
)

func NewServer(c Config) (*Server, error) {
	wire.Build(
		repository.New,
		dbmanager.NewMariaDBManager,
		dbmanager.NewMongoDBManager,
		repository.NewApplicationRepository,
		repository.NewGitRepositoryRepository,
		repository.NewEnvironmentRepository,
		repository.NewBuildRepository,
		repository.NewArtifactRepository,
		repository.NewUserRepository,
		grpc.NewAPIServiceServer,
		grpc.NewAuthInterceptor,
		grpc.NewControllerServiceClient,
		apiserver.NewService,
		provideRepositoryPublicKey,
		provideStorage,
		provideContainerLogger,
		provideGatewayServer,
		wire.FieldsOf(new(Config), "AvatarBaseURL", "AuthHeader", "Controller", "DB", "MariaDB", "MongoDB", "Storage"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
