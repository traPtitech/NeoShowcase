package discovery

import (
	"context"
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/util/hash"
)

type Cluster struct {
	d Discoverer

	targets []Target
	meIdx   int

	setInitialized func()
	initialized    chan struct{}
	lock           sync.RWMutex
}

func NewCluster(d Discoverer) *Cluster {
	initialized := make(chan struct{})
	return &Cluster{
		d:     d,
		meIdx: -1,

		setInitialized: sync.OnceFunc(func() {
			close(initialized)
		}),
		initialized: initialized,
	}
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
		c.setInitialized()
		log.Infof("[cluster] %d targets received", len(targets))
	}
	return nil
}

func (c *Cluster) IsLeader() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	<-c.initialized

	return c.meIdx == 0
}

func (c *Cluster) Me() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	<-c.initialized

	return c.meIdx
}

func (c *Cluster) Key(key string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	<-c.initialized

	return hash.JumpHashStr(key, len(c.targets))
}

func (c *Cluster) Assigned(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	<-c.initialized

	if len(c.targets) == 0 {
		return false
	}
	return c.Key(key) == c.meIdx
}

func (c *Cluster) AllNeighbors() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	<-c.initialized

	ips := make([]string, 0, len(c.targets)-1)
	for _, t := range c.targets {
		if t.Me {
			continue
		}
		ips = append(ips, t.IP)
	}
	return ips
}
