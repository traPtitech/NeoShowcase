// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

// Injectors from wire.go:

func New(c2 Config) (*Server, error) {
	config := provideAdminDBConfig(c2)
	db, err := admindb.New(config)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	storageConfig := provideStorageConfig(c2)
	storage, err := initStorage(storageConfig)
	if err != nil {
		return nil, err
	}
	webServerDocumentRootPath := provideWebServerDocumentRootPath(c2)
	webServerPort := provideWebServerPort(c2)
	engine := staticserver.NewBuiltIn(storage, webServerDocumentRootPath, webServerPort)
	buildRepository := repository.NewBuildRepository(db)
	staticSiteServerService := usecase.NewStaticSiteServerService(buildRepository, engine, db)
	staticSiteService := grpc.NewStaticSiteServiceServer(staticSiteServerService)
	tcpListenPort := provideGRPCPort(c2)
	mainServer := &Server{
		db:         db,
		grpcServer: server,
		engine:     engine,
		sss:        staticSiteService,
		port:       tcpListenPort,
	}
	return mainServer, nil
}
