package repository

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

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

func (u *userRepository) CreateUser(ctx context.Context, args domain.CreateUserArgs) (*domain.User, error) {
	user := &models.User{
		ID:   domain.NewID(),
		Name: args.Name,
	}
	if err := user.Insert(ctx, u.db, boil.Infer()); err != nil {
		return nil, errors.Wrap(err, "failed to insert user")
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := models.Users(models.UserWhere.ID.EQ(id)).One(ctx, u.db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
