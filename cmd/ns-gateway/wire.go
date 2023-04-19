//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func NewServer(c Config) (*Server, error) {
	wire.Build(
		admindb.New,
		dbmanager.NewMariaDBManager,
		dbmanager.NewMongoDBManager,
		repository.NewApplicationRepository,
		repository.NewAvailableDomainRepository,
		repository.NewGitRepositoryRepository,
		repository.NewEnvironmentRepository,
		repository.NewBuildRepository,
		repository.NewArtifactRepository,
		repository.NewUserRepository,
		grpc.NewAPIServiceServer,
		grpc.NewAuthInterceptor,
		grpc.NewControllerServiceClient,
		usecase.NewAPIServerService,
		provideRepositoryPublicKey,
		provideStorage,
		provideContainerLogger,
		provideGatewayServer,
		wire.FieldsOf(new(Config), "Controller", "DB", "MariaDB", "MongoDB", "Storage"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
