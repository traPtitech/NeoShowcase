package grpc

import (
	"context"
	"time"

	"connectrpc.com/connect/v2"
	"github.com/motoki317/sc"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type AuthInterceptor struct {
	userCache *sc.Cache[string, *domain.User]
	header    string
}

type AuthHeader string

func NewAuthInterceptor(
	userRepo domain.UserRepository,
	header AuthHeader,
) *AuthInterceptor {
	return &AuthInterceptor{
		userCache: sc.NewMust(func(ctx context.Context, name string) (*domain.User, error) {
			return userRepo.EnsureUser(ctx, name)
		}, 1*time.Minute, 2*time.Minute),
		header: string(header),
	}
}

func (a *AuthInterceptor) authenticate(ctx *context.Context) error {
	info, _ := connect.CallInfoForServerContext(*ctx)
	name := info.RequestHeader().Get(a.header)
	if name == "" {
		return connect.NewError(connect.CodeUnauthenticated, "")
	}
	user, err := a.userCache.Get(*ctx, name)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
	}
	web.SetUser(ctx, user)
	return nil
}

// Server returns the interceptor that resolves the user from the auth header.
func (a *AuthInterceptor) Server() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			err := a.authenticate(&ctx)
			if err != nil {
				return err
			}
			return next(ctx, spec, stream)
		}
	}
}
