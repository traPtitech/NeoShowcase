package retry

import (
	"context"
	"time"
)

func Do(ctx context.Context, fn func(ctx context.Context) error, initialBackoff, maxBackoff time.Duration) {
	backoff := initialBackoff / 2
	for {
		err := fn(ctx)
		if err == nil {
			backoff = initialBackoff
		}
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return
		}
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}
