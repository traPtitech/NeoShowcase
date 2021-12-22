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

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
type GitRepositoryRepository interface {
	RegisterRepository(ctx context.Context, args RegisterRepositoryArgs) (*domain.Repository, error)
	GetRepositoryByID(ctx context.Context, id string) (*domain.Repository, error)
	GetRepository(ctx context.Context, args GetRepositoryArgs) (*domain.Repository, error)
	RegisterProvider(ctx context.Context, args RegisterProviderArgs) (*domain.Provider, error)
	GetProviderByID(ctx context.Context, id string) (*domain.Provider, error)
	GetProviderByHost(ctx context.Context, host string) (*domain.Provider, error)
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

type RegisterProviderArgs struct {
	Domain string
	Secret string
}

type GetRepositoryArgs struct {
	ProviderID string
	Owner      string
	Name       string
}

func NewGitRepositoryRepository(db *sql.DB) GitRepositoryRepository {
	return &gitrepositoryRepository{
		db: db,
	}
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
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(repo.ProviderID)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
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

func (r *gitrepositoryRepository) GetRepository(ctx context.Context, args GetRepositoryArgs) (*domain.Repository, error) {
	const errMsg = "failed to GetRepository: %w"

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(args.ProviderID)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	repo, err := models.Repositories(models.RepositoryWhere.ProviderID.EQ(args.ProviderID), models.RepositoryWhere.Owner.EQ(args.Owner), models.RepositoryWhere.Name.EQ(args.Name)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
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

func (r *gitrepositoryRepository) RegisterProvider(ctx context.Context, args RegisterProviderArgs) (*domain.Provider, error) {
	const errMsg = "failed to RegisterProvider: %w"

	prov, err := models.Providers(models.ProviderWhere.Domain.EQ(args.Domain)).One(ctx, r.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if prov != nil {
		return nil, fmt.Errorf(errMsg, errors.New("provider already exists"))
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		prov = &models.Provider{
			ID:     id.String(),
			Domain: args.Domain,
			Secret: args.Secret,
		}
		if err := prov.Insert(ctx, r.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}

	}

	log.WithField("providerID", prov.ID).
		WithField("domain", prov.Domain).
		Info("registered provider")

	return &domain.Provider{
		ID:     prov.ID,
		Secret: prov.Secret,
	}, nil
}

func (r *gitrepositoryRepository) GetProviderByID(ctx context.Context, id string) (*domain.Provider, error) {
	const errMsg = "failed to GetProviderByID: %w"

	prov, err := models.Providers(models.ProviderWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.Provider{
		ID:     prov.ID,
		Secret: prov.Secret,
	}, nil
}

func (r *gitrepositoryRepository) GetProviderByHost(ctx context.Context, host string) (*domain.Provider, error) {
	const errMsg = "failed to GetProviderByHost: %w"

	prov, err := models.Providers(models.ProviderWhere.Domain.EQ(host)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf(errMsg, err)
	}

	return &domain.Provider{
		ID:     prov.ID,
		Secret: prov.Secret,
	}, nil
}
