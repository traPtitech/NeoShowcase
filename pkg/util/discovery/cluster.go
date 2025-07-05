package discovery

import (
	"context"
	"sync"

	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type Cluster struct {
	d Discoverer

	targets []Target
	meIdx   int

	lock sync.RWMutex
}

func NewCluster(d Discoverer) *Cluster {
	return &Cluster{d: d, meIdx: -1}
}

func (c *Cluster) Start(ctx context.Context) error {
	updates, err := c.d.Watch(ctx)
	if err != nil {
		return err
	}
	for targets := range updates {
		c.lock.Lock()
		c.targets = targets
		_, c.meIdx, _ = lo.FindIndexOf(targets, func(e Target) bool { return e.Me })
		c.lock.Unlock()
	}
	return nil
}

func (c *Cluster) IsLeader() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.meIdx == 0
}

func (c *Cluster) Assigned(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if len(c.targets) == 0 {
		return false
	}
	return hash.JumpHashStr(key, len(c.targets)) == c.meIdx
}

func (c *Cluster) AllNeighbors() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ips := make([]string, 0, len(c.targets)-1)
	for _, t := range c.targets {
		if t.Me {
			continue
		}
		ips = append(ips, t.IP)
	}
	return ips
}
