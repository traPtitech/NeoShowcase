package usecase

import (
	"context"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/motoki317/sc"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util/optional"
)

type ContainerStateMutator struct {
	appRepo domain.ApplicationRepository
	backend domain.Backend

	lock sync.Mutex
}

func NewContainerStateMutator(
	appRepo domain.ApplicationRepository,
	backend domain.Backend,
) *ContainerStateMutator {
	m := &ContainerStateMutator{
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
		appID := e.ApplicationID
		// coalesce events
		go func() {
			_, err := updateOne.Get(context.Background(), appID)
			if err != nil {
				log.Errorf("failed to update app container state: %+v", err)
			}
		}()
	}
}

func (m *ContainerStateMutator) _updateOne(ctx context.Context, appID string) (struct{}, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	state, err := m.backend.GetContainer(ctx, appID)
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to get container state")
	}
	err = m.appRepo.UpdateApplication(ctx, appID, &domain.UpdateApplicationArgs{Container: optional.From(state.State)})
	if err != nil {
		return struct{}{}, errors.Wrap(err, "failed to update application")
	}

	return struct{}{}, nil
}

func (m *ContainerStateMutator) updateAll(ctx context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	apps, err := m.appRepo.GetApplications(ctx, domain.GetApplicationCondition{DeployType: optional.From(domain.DeployTypeRuntime)})
	if err != nil {
		return errors.Wrap(err, "failed to get all runtime applications")
	}

	containers, err := m.backend.ListContainers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}
	stateMap := lo.SliceToMap(containers, func(c *domain.Container) (string, domain.ContainerState) {
		return c.ApplicationID, c.State
	})

	appStates := lo.SliceToMap(apps, func(app *domain.Application) (string, domain.ContainerState) {
		return app.ID, lo.ValueOr(stateMap, app.ID, domain.ContainerStateMissing)
	})
	err = m.appRepo.BulkUpdateState(ctx, appStates)
	if err != nil {
		return errors.Wrap(err, "failed to bulk update state")
	}
	return nil
}
