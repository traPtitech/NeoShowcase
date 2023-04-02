package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
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
		mu = fromDomainUser(user)
		err = mu.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert user")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to commit")
	}

	return toDomainUser(mu), nil
}
