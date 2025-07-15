package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
)

type repositoryCommitRepository struct {
	db *sql.DB
}

func NewRepositoryCommitRepository(db *sql.DB) domain.RepositoryCommitRepository {
	return &repositoryCommitRepository{db: db}
}

func (r *repositoryCommitRepository) GetCommits(ctx context.Context, hashes []string) ([]*domain.RepositoryCommit, error) {
	if len(hashes) == 0 {
		return nil, nil
	}

	commits, err := models.RepositoryCommits(
		models.RepositoryCommitWhere.Hash.IN(hashes),
	).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query repository commits")
	}

	return ds.Map(commits, repoconvert.ToDomainRepositoryCommit), nil
}

func (r *repositoryCommitRepository) RecordCommit(ctx context.Context, commit *domain.RepositoryCommit) error {
	m := repoconvert.FromDomainRepositoryCommit(commit)
	err := m.Insert(ctx, r.db, boil.Blacklist())
	if err != nil {
		return errors.Wrap(err, "failed to insert repository commit")
	}
	return nil
}
