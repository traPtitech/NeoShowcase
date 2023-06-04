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
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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

	if cond.IDs.Valid {
		mods = append(mods, models.RepositoryWhere.ID.IN(cond.IDs.V))
	}
	if cond.URLs.Valid {
		mods = append(mods, models.RepositoryWhere.URL.IN(cond.URLs.V))
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

	modelRepos, err := models.Repositories(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repositories")
	}

	repos := ds.Map(modelRepos, repoconvert.ToDomainRepository)

	if cond.PublicOrOwnedBy.Valid {
		userID := cond.PublicOrOwnedBy.V
		repos = lo.Filter(repos, func(repo *domain.Repository, _ int) bool {
			return lo.Contains(repo.OwnerIDs, userID) || !repo.Auth.Valid
		})
	}

	return repos, nil
}

func (r *gitRepositoryRepository) getRepository(ctx context.Context, ex boil.ContextExecutor, id string) (*models.Repository, error) {
	return models.Repositories(
		models.RepositoryWhere.ID.EQ(id),
		qm.Load(models.RepositoryRels.RepositoryAuth),
		qm.Load(models.RepositoryRels.Users),
	).One(ctx, ex)
}

func (r *gitRepositoryRepository) GetRepository(ctx context.Context, id string) (*domain.Repository, error) {
	repo, err := r.getRepository(ctx, r.db, id)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to get repository")
	}
	return repoconvert.ToDomainRepository(repo), nil
}

func (r *gitRepositoryRepository) CreateRepository(ctx context.Context, repo *domain.Repository) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer tx.Rollback()

	mr := repoconvert.FromDomainRepository(repo)
	err = mr.Insert(ctx, tx, boil.Blacklist())
	if err != nil {
		return errors.Wrap(err, "failed to insert repository")
	}

	if repo.Auth.Valid {
		mra := repoconvert.FromDomainRepositoryAuth(repo.ID, &repo.Auth.V)
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

func (r *gitRepositoryRepository) UpdateRepository(ctx context.Context, id string, args *domain.UpdateRepositoryArgs) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	repo, err := r.getRepository(ctx, tx, id)
	var cols []string

	if args.Name.Valid {
		repo.Name = args.Name.V
		cols = append(cols, models.RepositoryColumns.Name)
	}
	if args.URL.Valid {
		repo.URL = args.URL.V
		cols = append(cols, models.RepositoryColumns.URL)
	}

	if len(cols) > 0 {
		_, err = repo.Update(ctx, tx, boil.Whitelist(cols...))
		if err != nil {
			return errors.Wrap(err, "failed to update repository")
		}
	}

	if args.Auth.Valid {
		if repo.R != nil && repo.R.RepositoryAuth != nil {
			_, err = repo.R.RepositoryAuth.Delete(ctx, tx)
			if err != nil {
				return errors.Wrap(err, "failed to delete existing repository auth")
			}
		}
		if args.Auth.V.Valid {
			mra := repoconvert.FromDomainRepositoryAuth(repo.ID, &args.Auth.V.V)
			err = repo.SetRepositoryAuth(ctx, tx, true, mra)
			if err != nil {
				return errors.Wrap(err, "failed to set repository auth")
			}
		}
	}

	if args.OwnerIDs.Valid {
		err = r.setOwners(ctx, tx, repo, args.OwnerIDs.V)
		if err != nil {
			return err
		}
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

func (r *gitRepositoryRepository) DeleteRepository(ctx context.Context, id string) error {
	repo, err := r.getRepository(ctx, r.db, id)
	if err != nil {
		return err
	}
	err = repo.SetUsers(ctx, r.db, false)
	if err != nil {
		return errors.Wrap(err, "failed to delete repository owners")
	}
	if repo.R.RepositoryAuth != nil {
		_, err = repo.R.RepositoryAuth.Delete(ctx, r.db)
		if err != nil {
			return errors.Wrap(err, "failed to delete repository auth")
		}
	}
	_, err = repo.Delete(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "failed to delete repository")
	}
	return nil
}
