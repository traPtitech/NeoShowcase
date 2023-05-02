package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/traPtitech/neoshowcase/pkg/interface/repository/repoconvert"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUsers(ctx context.Context, cond domain.GetUserCondition) ([]*domain.User, error) {
	var mods []qm.QueryMod

	if cond.Admin.Valid {
		mods = append(mods, models.UserWhere.Admin.EQ(cond.Admin.V))
	}

	users, err := models.Users(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "getting users")
	}
	return lo.Map(users, func(u *models.User, _ int) *domain.User {
		return repoconvert.ToDomainUser(u)
	}), nil
}

func (r *userRepository) GetOrCreateUser(ctx context.Context, name string) (*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	mu, err := models.Users(
		models.UserWhere.Name.EQ(name),
		qm.For("UPDATE"),
	).One(ctx, tx)
	if err != nil && !isNoRowsErr(err) {
		return nil, errors.Wrap(err, "failed to get user")
	}

	if isNoRowsErr(err) {
		user := domain.NewUser(name)
		mu = repoconvert.FromDomainUser(user)
		err = mu.Insert(ctx, tx, boil.Blacklist())
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert user")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to commit")
	}

	return repoconvert.ToDomainUser(mu), nil
}

func (r *userRepository) GetUserKeys(ctx context.Context, cond domain.GetUserKeyCondition) ([]*domain.UserKey, error) {
	var mods []qm.QueryMod

	if cond.UserIDs.Valid {
		mods = append(mods, models.UserKeyWhere.UserID.IN(cond.UserIDs.V))
	}

	keys, err := models.UserKeys(mods...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "getting user keys")
	}
	return lo.Map(keys, func(key *models.UserKey, _ int) *domain.UserKey {
		return repoconvert.ToDomainUserKey(key)
	}), nil
}

func (r *userRepository) CreateUserKey(ctx context.Context, key *domain.UserKey) error {
	mk := repoconvert.FromDomainUserKey(key)
	err := mk.Insert(ctx, r.db, boil.Blacklist())
	if err != nil {
		return errors.Wrap(err, "inserting user key")
	}
	return nil
}

func (r *userRepository) DeleteUserKey(ctx context.Context, keyID string, userID string) error {
	_, err := models.UserKeys(
		models.UserKeyWhere.ID.EQ(keyID),
		models.UserKeyWhere.UserID.EQ(userID),
	).DeleteAll(ctx, r.db)
	if err != nil {
		return errors.Wrap(err, "deleting user key")
	}
	return nil
}
