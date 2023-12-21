package grpc

import (
	"connectrpc.com/connect"
	"context"
	"github.com/friendsofgo/errors"
	"net/http"
)

type TokenAuthInterceptor struct {
	header string
	token  string
}

var _ connect.Interceptor = &TokenAuthInterceptor{}

func NewTokenAuthInterceptor(
	header string,
	token string,
) (*TokenAuthInterceptor, error) {
	if header == "" {
		return nil, errors.New("header name cannot be empty")
	}
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}
	return &TokenAuthInterceptor{
		header: header,
		token:  token,
	}, nil
}

func (a *TokenAuthInterceptor) addHeader(header http.Header) {
	header.Add(a.header, a.token)
}

func (a *TokenAuthInterceptor) authenticate(headers http.Header) error {
	token := headers.Get(a.header)
	if token == "" || token != a.token {
		return connect.NewError(connect.CodeUnauthenticated, nil)
	}
	return nil
}

func (a *TokenAuthInterceptor) WrapUnary(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		if request.Spec().IsClient {
			a.addHeader(request.Header())
		} else {
			err := a.authenticate(request.Header())
			if err != nil {
				return nil, err
			}
		}
		return unaryFunc(ctx, request)
	}
}

func (a *TokenAuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)
		a.addHeader(conn.RequestHeader())
		return conn
	}
}

func (a *TokenAuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		err := a.authenticate(conn.RequestHeader())
		if err != nil {
			return err
		}
		return next(ctx, conn)
	}
}
