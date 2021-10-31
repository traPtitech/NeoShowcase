package repository

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, args CreateUserArgs) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type CreateUserArgs struct {
	Name string
}
