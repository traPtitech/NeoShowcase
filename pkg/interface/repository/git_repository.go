package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

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

func (r *gitRepositoryRepository) GetRepositories(ctx context.Context, cond domain.GetRepositoryCondition) ([]*domain.Repository, error) {
	mods := []qm.QueryMod{
		qm.Load(models.RepositoryRels.RepositoryAuth),
		qm.Load(models.RepositoryRels.Users),
	}

	if cond.UserID.Valid {
		mods = append(mods,
			qm.InnerJoin(fmt.Sprintf(
				"%s ON %s.repository_id = %s",
				models.TableNames.RepositoryOwners,
				models.TableNames.RepositoryOwners,
				models.RepositoryTableColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.user_id = ?", models.TableNames.RepositoryOwners), cond.UserID.V),
		)
	}

	repos, err := models.Repositories(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repositories")
	}
	return lo.Map(repos, func(repo *models.Repository, i int) *domain.Repository {
		return toDomainRepository(repo)
	}), nil
}

func (r *gitRepositoryRepository) GetRepository(ctx context.Context, id string) (*domain.Repository, error) {
	repo, err := models.Repositories(
		models.RepositoryWhere.ID.EQ(id),
		qm.Load(models.RepositoryRels.RepositoryAuth),
		qm.Load(models.RepositoryRels.Users),
	).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get repository")
	}
	return toDomainRepository(repo), nil
}

func (r *gitRepositoryRepository) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	mr := fromDomainRepository(repo)
	err = mr.Insert(ctx, tx, boil.Blacklist())
	if err != nil {
		return errors.Wrap(err, "failed to insert repository")
	}

	if repo.Auth.Valid {
		mra := fromDomainRepositoryAuth(repo.ID, &repo.Auth.V)
		err = mra.Insert(ctx, tx, boil.Blacklist())
		if err != nil {
			return errors.Wrap(err, "failed to insert repository auth")
		}
	}

	err = r.setOwners(ctx, tx, mr, repo.OwnerIDs)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *gitRepositoryRepository) setOwners(ctx context.Context, ex boil.ContextExecutor, repo *models.Repository, ownerIDs []string) error {
	ownerIDs = lo.Uniq(ownerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(ownerIDs)).All(ctx, ex)
	if err != nil {
		return errors.Wrap(err, "failed to get users")
	}
	if len(users) < len(ownerIDs) {
		return ErrNotFound
	}
	err = repo.SetUsers(ctx, ex, false, users...)
	if err != nil {
		return errors.Wrap(err, "failed to")
	}
	return nil
}
