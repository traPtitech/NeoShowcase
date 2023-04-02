package grpc

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/motoki317/sc"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type userKeyType struct{}

var userKey = userKeyType{}

func getUser(ctx context.Context) *domain.User {
	return ctx.Value(userKey).(*domain.User)
}

func setUser(ctx *context.Context, user *domain.User) {
	*ctx = context.WithValue(*ctx, userKey, user)
}

type AuthInterceptor struct {
	userCache *sc.Cache[string, *domain.User]
}

var _ connect.Interceptor = &AuthInterceptor{}

func NewAuthInterceptor(
	userRepo domain.UserRepository,
) *AuthInterceptor {
	return &AuthInterceptor{
		userCache: sc.NewMust(func(ctx context.Context, name string) (*domain.User, error) {
			return userRepo.GetOrCreateUser(ctx, name)
		}, 1*time.Minute, 2*time.Minute),
	}
}

func (a *AuthInterceptor) authenticate(ctx *context.Context, headers http.Header) error {
	name := headers.Get(web.HeaderNameAPIAuthorization)
	if name == "" {
		return connect.NewError(connect.CodeUnauthenticated, nil)
	}
	user, err := a.userCache.Get(*ctx, name)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	setUser(ctx, user)
	return nil
}

func (a *AuthInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		err := a.authenticate(&ctx, request.Header())
		if err != nil {
			return nil, err
		}
		return unaryFunc(ctx, request)
	}
}

func (a *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		err := a.authenticate(&ctx, conn.RequestHeader())
		if err != nil {
			return err
		}
		return next(ctx, conn)
	}
}
