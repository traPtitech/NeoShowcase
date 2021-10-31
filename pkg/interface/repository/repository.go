package repository

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type GitrepositoryRepository interface {
	CreateRepository(ctx context.Context, args *CreateRepositoryArgs) (*domain.Repository, error)
	GetRepositoryByID(ctx context.Context, id string) (*domain.Repository, error)
	GetRepositoryByOwnerAndName(ctx context.Context, owner, name string) (*domain.Repository, error)
	CreateProvider(ctx context.Context, args *CreateProviderArgs) (*domain.Provider, error)
	GetProviderByID(ctx context.Context, id string) (*domain.Provider, error)
	GetProvierByDomain(ctx context.Context, domain string) (*domain.Provider, error)
}

type CreateRepositoryArgs struct {
	Name       string
	Owner      string
	URL        string
	Secret     string
	ProviderID string // TODO: providerid型を作る
}

type CreateProviderArgs struct {
	Name   string
	Domain string
	Secret string // TODO: repoとproviderこれどっちだっけ？
}
