package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type gitRepositoryRepository struct {
	db *sql.DB
}

func NewGitRepositoryRepository(db *sql.DB) domain.GitRepositoryRepository {
	return &gitRepositoryRepository{
		db: db,
	}
}

func (r *gitRepositoryRepository) RegisterRepository(ctx context.Context, args domain.RegisterRepositoryArgs) (domain.Repository, error) {
	const errMsg = "failed to RegisterRepository: %w"

	repo, err := models.Repositories(models.RepositoryWhere.URL.EQ(args.URL)).One(ctx, r.db)
	if err != nil && err != sql.ErrNoRows {
		return domain.Repository{}, fmt.Errorf(errMsg, err)
	}
	if repo != nil {
		return domain.Repository{}, fmt.Errorf(errMsg, errors.New("repository already exists"))
	}

	repo = &models.Repository{
		ID:   domain.NewID(),
		Name: args.Name,
		URL:  args.URL,
	}
	if err := repo.Insert(ctx, r.db, boil.Infer()); err != nil {
		return domain.Repository{}, fmt.Errorf(errMsg, err)
	}

	log.WithField("repositoryID", repo.ID).
		Info("registered repository")

	return toDomainRepository(repo), nil

}

func (r *gitRepositoryRepository) GetRepositoryByID(ctx context.Context, id string) (domain.Repository, error) {
	const errMsg = "failed to GetRepositoryByID: %w"

	repo, err := models.Repositories(models.RepositoryWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return domain.Repository{}, ErrNotFound
		}
		return domain.Repository{}, fmt.Errorf(errMsg, err)
	}

	return toDomainRepository(repo), nil
}

func (r *gitRepositoryRepository) GetRepository(ctx context.Context, rawURL string) (domain.Repository, error) {
	const errMsg = "failed to GetRepository: %w"

	repo, err := models.Repositories(models.RepositoryWhere.URL.EQ(rawURL)).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return domain.Repository{}, ErrNotFound
		}
		return domain.Repository{}, fmt.Errorf(errMsg, err)
	}

	return toDomainRepository(repo), nil
}
