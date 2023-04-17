package loop

import (
	"context"
	"time"
)

func Loop(ctx context.Context, do func(ctx context.Context), interval time.Duration, first bool) {
	if first {
		do(ctx)
	}
	for {
		select {
		case <-time.After(interval):
			do(ctx)
		case <-ctx.Done():
			return
		}
	}
}
