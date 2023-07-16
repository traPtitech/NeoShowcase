package scutil

import (
	"context"
	"time"

	"github.com/motoki317/sc"
)

const inf = 100 * 365 * 24 * time.Hour

func Once[T any](fn func(ctx context.Context) (T, error)) func(ctx context.Context) (T, error) {
	// ref: sync.OnceFunc in go 1.21+
	replaceFn := func(ctx context.Context, _ struct{}) (T, error) {
		return fn(ctx)
	}
	c := sc.NewMust(replaceFn, inf, inf)
	return func(ctx context.Context) (T, error) {
		return c.Get(ctx, struct{}{})
	}
}
