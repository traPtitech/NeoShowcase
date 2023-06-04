//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

func NewServer(c Config) (*Server, error) {
	wire.Build(
		repository.New,
		repository.NewGitRepositoryRepository,
		repository.NewUserRepository,
		giteaintegration.NewIntegration,
		wire.FieldsOf(new(Config), "Gitea", "DB"),
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
