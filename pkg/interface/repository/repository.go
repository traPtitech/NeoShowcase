package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/apex/log"
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
	Secret string // TODO: repoとproviderこれどっちだっけ？
}

func (r *gitrepositoryRepository) RegisterRepository(ctx context.Context, args RegisterRepositoryArgs) (*domain.Repository, error) {
	const errMsg = "failed to RegisterRepository: %w"

	repo, err := models.Repositories(models.RepositoryWhere.URL.EQ(args.URL)).One(ctx, r.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if repo != nil {
		return nil, fmt.Errorf(errMsg, "repository already exists")
	} else {
		repo = &models.Repository{
			RepositoryID: uuid.NewRandom().String(),
			Owner:        args.RepositoryOwner,
			Name:         args.RepositoryName,
			URL:          args.URL,
			ProviderID:   args.ProviderID,
		}
		if err := repo.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	log.WithField("repositoryID", repo.RepositoryID).
		Info("registered repository")

	return &domain.Repository{
		RepositoryID: repo.RepositoryID,
		Owner:        repo.Owner,
		Name:         repo.Name,
		URL:          repo.URL,
		ProviderID:   repo.ProviderID,
	}, nil

}
