package scutil

import (
	"context"
	"time"

	"github.com/motoki317/sc"
)

func WrapFunc[T any](fn func(ctx context.Context) (T, error)) func(ctx context.Context, _ struct{}) (T, error) {
	return func(ctx context.Context, _ struct{}) (T, error) {
		return fn(ctx)
	}
}

func Wrap[T any](fn func(ctx context.Context) (T, error), freshFor, ttl time.Duration) func(ctx context.Context) (T, error) {
	replaceFn := func(ctx context.Context, _ struct{}) (T, error) {
		return fn(ctx)
	}
	c := sc.NewMust(replaceFn, freshFor, ttl)
	return func(ctx context.Context) (T, error) {
		return c.Get(ctx, struct{}{})
	}
}
