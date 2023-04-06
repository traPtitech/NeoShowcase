package coalesce

import (
	"context"

	"github.com/motoki317/sc"
)

type Coalescer struct {
	c *sc.Cache[struct{}, struct{}]
}

func (c *Coalescer) Do(ctx context.Context) error {
	_, err := c.c.Get(ctx, struct{}{})
	return err
}

func NewCoalescer(fn func(ctx context.Context) error) *Coalescer {
	wrappedFn := func(ctx context.Context, _ struct{}) (struct{}, error) {
		return struct{}{}, fn(ctx)
	}
	return &Coalescer{
		c: sc.NewMust(wrappedFn, 0, 0, sc.EnableStrictCoalescing()),
	}
}
