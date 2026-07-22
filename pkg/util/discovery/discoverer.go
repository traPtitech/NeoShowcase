package discovery

import (
	"context"

	"github.com/friendsofgo/errors"
)

type Target struct {
	IP string
	Me bool
}

type Discoverer interface {
	// Watch watches for changes in the list of IPs and sends updates through the returned channel.
	// The context can be used to stop watching.
	//
	// Returned results are sorted by IP and are expected to pass the validateTargets func.
	Watch(ctx context.Context) (<-chan []Target, error)
}

func validateTargets(targets []Target) error {
	meCnt := 0
	for _, t := range targets {
		if t.Me {
			meCnt++
		}
	}
	// Note that it is possible there are no targets marked as "me".
	if meCnt > 1 {
		return errors.Errorf("too many targets marked as \"me\": %d", meCnt)
	}
	for i := 0; i < len(targets)-1; i++ {
		if targets[i].IP == targets[i+1].IP {
			return errors.Errorf("duplicate target: %s", targets[i].IP)
		}
	}
	return nil
}
