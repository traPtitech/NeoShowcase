package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/models"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository/repoconvert"
	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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
	return ds.Map(users, repoconvert.ToDomainUser), nil
}

func (r *userRepository) EnsureUser(ctx context.Context, name string) (*domain.User, error) {
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

func (r *userRepository) EnsureUsers(ctx context.Context, names []string) ([]*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback()

	mus, err := models.Users(
		models.UserWhere.Name.IN(names),
		qm.For("UPDATE"),
	).All(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "getting users")
	}

	modelUsers := lo.SliceToMap(mus, func(mu *models.User) (string, *models.User) { return mu.Name, mu })
	for _, name := range names {
		if _, ok := modelUsers[name]; ok {
			continue
		}
		user := domain.NewUser(name)
		mu := repoconvert.FromDomainUser(user)
		err = mu.Insert(ctx, tx, boil.Blacklist())
		if err != nil {
			return nil, errors.Wrap(err, "inserting user")
		}
		modelUsers[name] = mu
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to commit")
	}

	return lo.MapToSlice(modelUsers, func(name string, mu *models.User) *domain.User {
		return repoconvert.ToDomainUser(mu)
	}), nil
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
	return ds.Map(keys, repoconvert.ToDomainUserKey), nil
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
