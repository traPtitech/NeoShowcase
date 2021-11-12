package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository interface {
	CreateUser(ctx context.Context, args CreateUserArgs) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

type CreateUserArgs struct {
	Name string
}

func (u *userRepository) CreateUser(ctx context.Context, args CreateUserArgs) (*domain.User, error) {
	const errMsg = "failed to create user: %w"

	user, err := models.Users(models.UserWhere.Name.EQ(args.Name)).One(ctx, u.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf(errMsg, err)
	} else if user != nil {
		return nil, fmt.Errorf(errMsg, errors.New("user already exists"))
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		user = &models.User{
			ID:   id.String(),
			Name: args.Name,
		}
		if err := user.Insert(ctx, u.db, boil.Infer()); err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	const errMsg = "failed to get user: %w"

	user, err := models.Users(models.UserWhere.ID.EQ(id)).One(ctx, u.db)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &domain.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}
