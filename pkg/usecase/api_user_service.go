package usecase

import (
	"context"

	"github.com/friendsofgo/errors"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

func (s *APIServerService) GetMe(ctx context.Context) *domain.User {
	return web.GetUser(ctx)
}

func (s *APIServerService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetUsers(ctx, domain.GetUserCondition{})
}

func (s *APIServerService) CreateUserKey(ctx context.Context, publicKey string) (*domain.UserKey, error) {
	user := web.GetUser(ctx)
	key, err := domain.NewUserKey(user.ID, publicKey)
	if err != nil {
		return nil, newError(ErrorTypeBadRequest, "invalid public key", err)
	}
	err = s.userRepo.CreateUserKey(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "creating user key")
	}
	return key, nil
}

func (s *APIServerService) GetUserKeys(ctx context.Context) ([]*domain.UserKey, error) {
	user := web.GetUser(ctx)
	keys, err := s.userRepo.GetUserKeys(ctx, domain.GetUserKeyCondition{
		UserIDs: optional.From([]string{user.ID}),
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting user keys")
	}
	return keys, nil
}

func (s *APIServerService) DeleteUserKey(ctx context.Context, keyID string) error {
	user := web.GetUser(ctx)
	return s.userRepo.DeleteUserKey(ctx, keyID, user.ID)
}
