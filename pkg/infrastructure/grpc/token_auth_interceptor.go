package grpc

import (
	"context"

	"connectrpc.com/connect/v2"
	"github.com/samber/oops"
)

type TokenAuthInterceptor struct {
	header string
	token  string
}

func NewTokenAuthInterceptor(
	header string,
	token string,
) (*TokenAuthInterceptor, error) {
	if header == "" {
		return nil, oops.New("header name cannot be empty")
	}
	if token == "" {
		return nil, oops.New("token cannot be empty")
	}
	return &TokenAuthInterceptor{
		header: header,
		token:  token,
	}, nil
}

// Client returns the interceptor that attaches the token to outgoing requests.
func (a *TokenAuthInterceptor) Client() connect.ClientInterceptor {
	return func(next connect.ClientFunc) connect.ClientFunc {
		return func(ctx context.Context, spec connect.Spec) (connect.ClientStream, error) {
			info, _ := connect.CallInfoForClientContext(ctx)
			info.RequestHeader().Add(a.header, a.token)
			return next(ctx, spec)
		}
	}
}

// Server returns the interceptor that rejects requests without the token.
func (a *TokenAuthInterceptor) Server() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			info, _ := connect.CallInfoForServerContext(ctx)
			token := info.RequestHeader().Get(a.header)
			if token == "" || token != a.token {
				return connect.NewError(connect.CodeUnauthenticated, "")
			}
			return next(ctx, spec, stream)
		}
	}
}
