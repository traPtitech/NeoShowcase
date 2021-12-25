//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leandro-lugaresi/hub"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/eventbus"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewServer,
		grpc.NewBuilderServiceServer,
		usecase.NewBuilderService,
		repository.NewBuildLogRepository,
		eventbus.NewLocal,
		admindb.New,
		hub.New,
		provideGRPCPort,
		provideDockerImageRegistry,
		provideStorageConfig,
		provideAdminDBConfig,
		initStorage,
		initBuildkitClient,
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
