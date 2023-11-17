package retry

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
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
			log.Infof("Retrier: retrying in %v: %v", backoff, msg)
		} else {
			log.Errorf("Retrier: retrying in %v: %v: %v", backoff, msg, err)
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
