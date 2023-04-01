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
		grpc.NewComponentServiceClient,
		usecase.NewStaticSiteServerService,
		admindb.New,
		repository.NewApplicationRepository,
		repository.NewBuildRepository,
		staticserver.NewBuiltIn,
		provideWebServerPort,
		provideWebServerDocumentRootPath,
		initStorage,
		wire.FieldsOf(new(Config), "NS", "DB", "Storage"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
