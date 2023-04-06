package dockerimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) syncAppContainer(ctx context.Context, app *domain.AppDesiredState, oldContainer *docker.APIContainers) error {
	newImageName := app.ImageName + ":" + app.ImageTag
	var oldRestartedAt time.Time
	var err error
	if oldContainer != nil {
		oldRestartedAt, err = time.Parse(time.RFC3339, oldContainer.Labels[appRestartedAtLabel])
		if err != nil {
			oldRestartedAt = time.Time{}
		}
	} else {
		oldRestartedAt = time.Time{}
	}
	doDeploy := oldContainer == nil || oldContainer.Image != newImageName || !oldRestartedAt.Equal(app.App.UpdatedAt)
	if !doDeploy {
		return nil
	}

	if oldContainer != nil {
		err = b.c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            oldContainer.ID,
			RemoveVolumes: true,
			Force:         true,
			Context:       ctx,
		})
		if err != nil {
			return errors.Wrap(err, "failed to remove old container")
		}
	}

	err = b.c.PullImage(docker.PullImageOptions{
		Repository: app.ImageName,
		Tag:        app.ImageTag,
		Context:    ctx,
	}, docker.AuthConfiguration{})
	if err != nil {
		return errors.Wrap(err, "failed to pull image")
	}

	envs := lo.MapToSlice(app.Envs, func(key string, value string) string {
		return key + "=" + value
	})
	cont, err := b.c.CreateContainer(docker.CreateContainerOptions{
		Name: containerName(app.App.ID),
		Config: &docker.Config{
			Image:  newImageName,
			Labels: containerLabels(app.App),
			Env:    envs,
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.RestartOnFailure(5),
		},
		NetworkingConfig: &docker.NetworkingConfig{EndpointsConfig: map[string]*docker.EndpointConfig{
			appNetwork: {
				Aliases: []string{networkName(app.App.ID)},
			},
		}},
		Context: ctx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	err = b.c.StartContainer(cont.ID, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start container")
	}
	return nil
}

func (b *dockerBackend) SynchronizeRuntime(ctx context.Context, apps []*domain.AppDesiredState) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	// List old resources
	oldContainers, err := b.c.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {fmt.Sprintf("%s=true", appLabel)}},
		Context: ctx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}
	oldContainersMap := lo.SliceToMap(oldContainers, func(c docker.APIContainers) (string, *docker.APIContainers) {
		return c.Labels[appIDLabel], &c
	})

	// Calculate next resources and apply
	for _, app := range apps {
		err = b.syncAppContainer(ctx, app, oldContainersMap[app.App.ID])
		if err != nil {
			log.Errorf("failed to sync app: %+v", err)
			continue // fail-safe
		}
	}

	// Synchronize ingress config
	cb := newRuntimeConfigBuilder()
	for _, app := range apps {
		for _, website := range app.App.Websites {
			cb.addWebsite(app.App, website)
		}
	}
	err = b.writeConfig(traefikRuntimeFilename, cb.build())
	if err != nil {
		return errors.Wrap(err, "failed to write runtime ingress config")
	}

	// Prune old resources
	newApps := lo.SliceToMap(apps, func(app *domain.AppDesiredState) (string, bool) { return app.App.ID, true })
	for _, oldContainer := range oldContainers {
		appID := oldContainer.Labels[appIDLabel]
		if newApps[appID] {
			continue
		}

		err = b.c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            oldContainer.ID,
			RemoveVolumes: true,
			Force:         true,
			Context:       ctx,
		})
		if err != nil {
			log.Errorf("failed to remove old container: %+v", err)
			continue // fail-safe
		}
	}

	return nil
}
