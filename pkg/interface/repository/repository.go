package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type GitrepositoryRepository interface {
	RegisterRepository(ctx context.Context, args *RegisterRepositoryArgs) (*domain.Repository, error)
	GetRepositoryByID(ctx context.Context, id string) (*domain.Repository, error)
	GetRepositoryByOwnerAndName(ctx context.Context, owner, name string) (*domain.Repository, error)
	CreateProvider(ctx context.Context, args *CreateProviderArgs) (*domain.Provider, error)
	GetProviderByID(ctx context.Context, id string) (*domain.Provider, error)
	GetProvierByDomain(ctx context.Context, domain string) (*domain.Provider, error)
}

type gitrepositoryRepository struct {
	db *sql.DB
}

type RegisterRepositoryArgs struct {
	RepositoryName  string
	RepositoryOwner string
	URL             string
	ProviderID      string // TODO: providerid型を作る
}

type CreateProviderArgs struct {
	Name   string
	Domain string
	Secret string
}

func (r *gitrepositoryRepository) RegisterRepository(ctx context.Context, args RegisterRepositoryArgs) (*domain.Repository, error) {
	const errMsg = "failed to RegisterRepository: %w"

	repo, err := models.Repositories(models.RepositoryWhere.URL.EQ(args.URL)).One(ctx, r.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if repo != nil {
		return nil, fmt.Errorf(errMsg, errors.New("repository already exists"))
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		repo = &models.Repository{
			ID:         id.String(),
			Owner:      args.RepositoryOwner,
			Name:       args.RepositoryName,
			URL:        args.URL,
			ProviderID: args.ProviderID,
		}
		if err := repo.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

	}

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(args.ProviderID)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	log.WithField("repositoryID", repo.ID).
		WithField("providerID", prov.ID).
		Info("registered repository")

	return &domain.Repository{
		ID:        repo.ID,
		RemoteURL: repo.URL,
		Provider: domain.Provider{
			ID:     prov.ID,
			Secret: prov.Secret,
		},
	}, nil

}

func (r *gitrepositoryRepository) GetRepositoryByID(ctx context.Context, id string) (*domain.Repository, error) {
	const errMsg = "failed to GetRepositoryByID: %w"

	repo, err := models.Repositories(models.RepositoryWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(repo.ProviderID)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &domain.Repository{
		ID:        repo.ID,
		RemoteURL: repo.URL,
		Provider: domain.Provider{
			ID:     prov.ID,
			Secret: prov.Secret,
		},
	}, nil
}

func (r *gitrepositoryRepository) GetRepositoryByOwnerAndName(ctx context.Context, owner, name string) (*domain.Repository, error) {
	const errMsg = "failed to GetRepositoryByOwnerAndName: %w"

	repo, err := models.Repositories(models.RepositoryWhere.Owner.EQ(owner), models.RepositoryWhere.Name.EQ(name)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(repo.ProviderID)).One(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &domain.Repository{
		ID:        repo.ID,
		RemoteURL: repo.URL,
		Provider: domain.Provider{
			ID:     prov.ID,
			Secret: prov.Secret,
		},
	}, nil
}
