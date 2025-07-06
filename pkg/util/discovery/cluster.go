package discovery

import (
	"context"
	"strconv"
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

func (c *Cluster) toAddr(ip string, port int) string {
	return "http://" + ip + ":" + strconv.Itoa(port)
}

func (c *Cluster) Start(ctx context.Context) error {
	updates, err := c.d.Watch(ctx)
	if err != nil {
		return err
	}
	for targets := range updates {
		c.lock.Lock()
		if len(targets) == 0 {
			log.Info("[cluster] no targets received, skipping until first discovery...")
			c.lock.Unlock()
			continue
		}
		c.targets = targets
		_, c.meIdx, _ = lo.FindIndexOf(targets, func(e Target) bool { return e.Me })
		c.lock.Unlock()
		c.setInitialized()
		log.Infof("[cluster] %d target(s) received", len(targets))
	}
	return nil
}

func (c *Cluster) IsLeader() bool {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.meIdx == 0
}

func (c *Cluster) Size() int {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.targets)
}

func (c *Cluster) Me() int {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.meIdx
}

func (c *Cluster) MyAddress(port int) (addr string, ok bool) {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.meIdx < 0 {
		return "", false
	}
	return c.toAddr(c.targets[c.meIdx].IP, port), true
}

func (c *Cluster) Key(key string) int {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	return hash.JumpHashStr(key, len(c.targets))
}

func (c *Cluster) Assigned(key string) bool {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	if len(c.targets) == 0 {
		return false
	}
	return c.Key(key) == c.meIdx
}

func (c *Cluster) AllNeighborAddresses(port int) []string {
	<-c.initialized
	c.lock.RLock()
	defer c.lock.RUnlock()

	ips := make([]string, 0, len(c.targets)-1)
	for _, t := range c.targets {
		if t.Me {
			continue
		}
		ips = append(ips, c.toAddr(t.IP, port))
	}
	return ips
}
