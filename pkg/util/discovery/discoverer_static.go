package discovery

import "context"

type staticDiscoverer struct {
	myIP string
}

func NewSingleDiscoverer(myIP string) Discoverer {
	return &staticDiscoverer{myIP: myIP}
}

func (d *staticDiscoverer) Watch(ctx context.Context) (<-chan []Target, error) {
	updates := make(chan []Target)
	targets := []Target{{IP: d.myIP, Me: true}}
	go func() {
		updates <- targets
		<-ctx.Done()
		close(updates)
	}()
	return updates, nil
}
