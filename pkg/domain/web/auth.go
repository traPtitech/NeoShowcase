package web

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type userKeyType struct{}

var userKey = userKeyType{}

func GetUser(ctx context.Context) *domain.User {
	return ctx.Value(userKey).(*domain.User)
}

func SetUser(ctx *context.Context, user *domain.User) {
	*ctx = context.WithValue(*ctx, userKey, user)
}
