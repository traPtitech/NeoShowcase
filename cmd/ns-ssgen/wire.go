//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

func New(c Config) (*Server, error) {
	wire.Build(
		grpc.NewServer,
		grpc.NewStaticSiteServiceServer,
		usecase.NewStaticSiteServerService,
		admindb.New,
		repository.NewBuildRepository,
		staticserver.NewBuiltIn,
		provideGRPCPort,
		provideStorageConfig,
		provideAdminDBConfig,
		provideWebServerPort,
		provideWebServerDocumentRootPath,
		initStorage,
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
