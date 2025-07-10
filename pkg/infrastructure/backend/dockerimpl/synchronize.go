package dockerimpl

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *Backend) Synchronize(ctx context.Context, s *domain.DesiredState) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	err := b.synchronizeRuntime(ctx, s.Runtime)
	if err != nil {
		return err
	}
	return b.synchronizeSSIngress(ctx, s.StaticSites)
}

func (b *Backend) SynchronizeShared(_ context.Context, _ *domain.DesiredStateLeader) error {
	// No shared resources (certificates) are required in docker traefik backend for now
	return nil
}
