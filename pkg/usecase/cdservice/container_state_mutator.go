package cdservice

import (
	"context"
	"log/slog"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/motoki317/sc"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/discovery"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ContainerStateMutator struct {
	cluster *discovery.Cluster
	appRepo domain.ApplicationRepository
	backend domain.Backend

	lock sync.Mutex
}

func NewContainerStateMutator(
	cluster *discovery.Cluster,
	appRepo domain.ApplicationRepository,
	backend domain.Backend,
) *ContainerStateMutator {
	m := &ContainerStateMutator{
		cluster: cluster,
		appRepo: appRepo,
		backend: backend,
	}
	go m._subscribe(backend)
	return m
}

func (m *ContainerStateMutator) _subscribe(backend domain.Backend) {
	sub, _ := backend.ListenContainerEvents()
	updateOne := sc.NewMust[string, struct{}](m._updateOne, 0, 0, sc.EnableStrictCoalescing())
	for e := range sub {
		// coalesce events
		go func(appID string) {
			_, err := updateOne.Get(context.Background(), appID)
			if err != nil {
				slog.Error("failed to update app container state", "error", err)
			}
		}(e.ApplicationID)
	}
}

func (m *ContainerStateMutator) _updateOne(ctx context.Context, appID string) (struct{}, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	container, err := m.backend.GetContainer(ctx, appID)
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to get container state")
	}
	err = m.appRepo.UpdateApplication(ctx, appID, &domain.UpdateApplicationArgs{
		Container:        optional.From(container.State),
		ContainerMessage: optional.From(container.Message),
	})
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to update application")
	}

	return struct{}{}, nil
}

func (m *ContainerStateMutator) updateAll(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Fetch all runtime apps
	allRuntimeApps, err := m.appRepo.GetApplications(ctx, domain.GetApplicationCondition{DeployType: optional.From(domain.DeployTypeRuntime)})
	if err != nil {
		return errors.Wrap(err, "failed to get all runtime applications")
	}
	// Shard by app ID
	allRuntimeApps = lo.Filter(allRuntimeApps, func(app *domain.Application, _ int) bool {
		return m.cluster.IsAssigned(app.ID)
	})

	// Fetch actual states
	containers, err := m.backend.ListContainers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}

	// If actual state is not found, update state as "missing"
	stateExists := lo.SliceToMap(containers, func(c *domain.Container) (string, bool) {
		return c.ApplicationID, true
	})
	for _, app := range allRuntimeApps {
		if !stateExists[app.ID] {
			containers = append(containers, &domain.Container{
				ApplicationID: app.ID,
				State:         domain.ContainerStateMissing,
			})
		}
	}

	// Update
	err = m.appRepo.BulkUpdateState(ctx, containers)
	if err != nil {
		return errors.Wrap(err, "failed to bulk update state")
	}
	return nil
}
