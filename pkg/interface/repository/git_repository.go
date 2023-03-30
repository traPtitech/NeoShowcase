package repository

import (
	"context"
	"database/sql"
	"fmt"

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
				"%s ON %s.%s = %s.%s",
				models.TableNames.RepositoryOwners,
				models.TableNames.RepositoryOwners,
				"repository_id",
				models.TableNames.Repositories,
				models.RepositoryColumns.ID,
			)),
			qm.Where(fmt.Sprintf("%s.%s = ?", models.TableNames.RepositoryOwners, "user_id"), cond.UserID.V),
		)
	}

	repos, err := models.Repositories(mods...).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
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
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return toDomainRepository(repo), nil
}

func (r *gitRepositoryRepository) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	mr := fromDomainRepository(repo)
	err = mr.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("faield to insert repository: %w", err)
	}

	if repo.Auth.Valid {
		mra := fromDomainRepositoryAuth(repo.ID, &repo.Auth.V)
		err = mr.SetRepositoryAuth(ctx, tx, true, mra)
		if err != nil {
			return fmt.Errorf("failed to insert repository auth: %w", err)
		}
	}

	repo.OwnerIDs = lo.Uniq(repo.OwnerIDs)
	users, err := models.Users(models.UserWhere.ID.IN(repo.OwnerIDs)).All(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	if len(users) < len(repo.OwnerIDs) {
		return ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *gitRepositoryRepository) RegisterRepositoryOwner(ctx context.Context, repositoryID string, userID string) error {
	repo, err := models.Repositories(models.RepositoryWhere.ID.EQ(repositoryID)).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get repository: %w", err)
	}
	user, err := models.Users(models.UserWhere.ID.EQ(userID)).One(ctx, r.db)
	if err != nil {
		if isNoRowsErr(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	err = repo.AddUsers(ctx, r.db, false, user)
	if err != nil {
		return fmt.Errorf("failed to add owner: %w", err)
	}
	return nil
}
