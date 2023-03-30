package retry

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func Do(ctx context.Context, fn func(ctx context.Context) error, initialBackoff, maxBackoff time.Duration) {
	backoff := initialBackoff
	for {
		err := fn(ctx)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			backoff = initialBackoff
			log.Infof("Lost connection, retrying in %v", backoff)
		} else {
			log.WithError(err).Errorf("Lost connection, retrying in %v", backoff)
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
