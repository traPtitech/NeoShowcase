package dockerimpl

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *dockerBackend) syncAppContainer(ctx context.Context, app *domain.RuntimeDesiredState, oldContainer *types.Container) error {
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
		err = b.c.ContainerRemove(ctx, oldContainer.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			return errors.Wrap(err, "failed to remove old container")
		}
	}

	registryAuth, err := b.authConfig()
	if err != nil {
		return errors.Wrap(err, "getting auth config")
	}
	res, err := b.c.ImagePull(ctx, app.ImageName+":"+app.ImageTag, types.ImagePullOptions{
		RegistryAuth: registryAuth,
	})
	if err != nil {
		return errors.Wrap(err, "pulling image")
	}
	_, err = io.ReadAll(res)
	if err != nil {
		return errors.Wrap(err, "pulling image")
	}
	err = res.Close()
	if err != nil {
		return errors.Wrap(err, "pulling image")
	}

	envs := lo.MapToSlice(app.Envs, func(key string, value string) string {
		return key + "=" + value
	})
	config := &container.Config{
		Image:  newImageName,
		Labels: b.containerLabels(app.App),
		Env:    envs,
	}
	if app.App.Config.Entrypoint != "" {
		config.Entrypoint = app.App.Config.EntrypointArgs()
	}
	if app.App.Config.Command != "" {
		config.Cmd = app.App.Config.CommandArgs()
	}
	hostConfig := &container.HostConfig{
		PortBindings: make(map[nat.Port][]nat.PortBinding),
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 5,
		},
	}
	for _, p := range app.App.PortPublications {
		appPort := nat.Port(fmt.Sprintf("%d/%s", p.ApplicationPort, p.Protocol))
		hostConfig.PortBindings[appPort] = append(hostConfig.PortBindings[appPort], nat.PortBinding{
			HostIP:   "0.0.0.0/0",
			HostPort: strconv.Itoa(p.InternetPort),
		})
	}
	if b.config.Resources.CPUs != 0 {
		hostConfig.NanoCPUs = int64(b.config.Resources.CPUs * 1e9)
	}
	if b.config.Resources.Memory != 0 {
		hostConfig.Memory = b.config.Resources.Memory
	}
	if b.config.Resources.MemorySwap != 0 {
		hostConfig.MemorySwap = b.config.Resources.MemorySwap
	}
	if b.config.Resources.MemoryReservation != 0 {
		hostConfig.MemoryReservation = b.config.Resources.MemoryReservation
	}
	networkingConfig := &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{
		b.config.Network: {
			Aliases: []string{networkName(app.App.ID)},
		},
	}}
	cont, err := b.c.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName(app.App.ID))
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	err = b.c.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start container")
	}
	return nil
}

func (b *dockerBackend) synchronizeRuntime(ctx context.Context, apps []*domain.RuntimeDesiredState) error {
	// List old resources
	oldContainers, err := b.c.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("label", fmt.Sprintf("%s=true", appLabel)),
		),
	})
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}
	oldContainersMap := lo.SliceToMap(oldContainers, func(c types.Container) (string, *types.Container) {
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
			cb.addWebsite(b, app.App, website)
		}
	}
	err = b.writeConfig(traefikRuntimeFilename, cb.build())
	if err != nil {
		return errors.Wrap(err, "failed to write runtime ingress config")
	}

	// Prune old resources
	newApps := lo.SliceToMap(apps, func(app *domain.RuntimeDesiredState) (string, bool) { return app.App.ID, true })
	for _, oldContainer := range oldContainers {
		appID := oldContainer.Labels[appIDLabel]
		if newApps[appID] {
			continue
		}

		err = b.c.ContainerRemove(ctx, oldContainer.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			log.Errorf("failed to remove old container: %+v", err)
			continue // fail-safe
		}
	}

	return nil
}
