package dockerimpl

import (
	"context"
	"fmt"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) Synchronize(ctx context.Context, apps []*domain.AppDesiredState) error {
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

	// Calculate next resources to apply
	for _, app := range apps {
		oldContainer, oldExist := oldContainersMap[app.App.ID]
		newImageName := app.ImageName + ":" + app.ImageTag
		doDeploy := !oldExist || app.Restart || oldContainer.Image != newImageName
		if !doDeploy {
			continue
		}

		if oldExist {
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
				Labels: containerLabels(app.App.ID),
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

		if err := b.c.StartContainer(cont.ID, nil); err != nil {
			return errors.Wrap(err, "failed to start container")
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
			return errors.Wrap(err, "failed to remove old container")
		}
	}

	return nil
}
