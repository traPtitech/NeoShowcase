package retry

import (
	"context"
	"log/slog"
	"time"
)

const (
	initialBackoff   = 1 * time.Second
	maxBackoff       = 60 * time.Second
	successThreshold = 60 * time.Second
)

func Do(ctx context.Context, fn func(ctx context.Context) error, msg string) {
	backoff := initialBackoff
	for {
		start := time.Now()
		err := fn(ctx)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if time.Since(start) >= successThreshold || err == nil {
			backoff = initialBackoff
		}
		if err == nil {
			slog.InfoContext(ctx, "Retrier: retrying", "backoff", backoff, "message", msg)
		} else {
			slog.ErrorContext(ctx, "Retrier: retrying", "backoff", backoff, "message", msg, "error", err)
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
