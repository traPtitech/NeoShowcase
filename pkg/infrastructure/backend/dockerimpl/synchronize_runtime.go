package dockerimpl

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/netip"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func (b *Backend) syncAppContainer(ctx context.Context, app *domain.RuntimeDesiredState, oldContainer *container.Summary) error {
	newImageName := app.ImageName + ":" + app.ImageTag
	oldRestartedAt := getRestartedAt(oldContainer)
	doDeploy := oldContainer == nil || oldContainer.Image != newImageName || !oldRestartedAt.Equal(app.App.UpdatedAt)
	if !doDeploy {
		return nil
	}

	if oldContainer != nil {
		_, err := b.c.ContainerRemove(ctx, oldContainer.ID, client.ContainerRemoveOptions{
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
	res, err := b.c.ImagePull(ctx, app.ImageName+":"+app.ImageTag, client.ImagePullOptions{
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
		Image:        newImageName,
		Labels:       b.containerLabels(app.App),
		Env:          envs,
		ExposedPorts: make(network.PortSet),
		OpenStdin:    true,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}
	if args, _ := domain.ParseArgs(app.App.Config.BuildConfig.GetRuntimeConfig().Entrypoint); len(args) > 0 {
		config.Entrypoint = args
	}
	if args, _ := domain.ParseArgs(app.App.Config.BuildConfig.GetRuntimeConfig().Command); len(args) > 0 {
		config.Cmd = args
	}
	for _, website := range app.App.Websites {
		config.ExposedPorts[network.MustParsePort(fmt.Sprintf("%d/tcp", website.HTTPPort))] = struct{}{}
	}
	for _, p := range app.App.PortPublications {
		config.ExposedPorts[network.MustParsePort(fmt.Sprintf("%d/%s", p.ApplicationPort, p.Protocol))] = struct{}{}
	}
	hostConfig := &container.HostConfig{
		PortBindings: make(network.PortMap),
		RestartPolicy: container.RestartPolicy{
			Name: "on-failure",
			// sablier stops the container, so we don't need to restart it
			MaximumRetryCount: lo.Ternary(b.useSablier(app.App), 0, 5),
		},
	}
	for _, p := range app.App.PortPublications {
		appPort := network.MustParsePort(fmt.Sprintf("%d/%s", p.ApplicationPort, p.Protocol))
		hostConfig.PortBindings[appPort] = append(hostConfig.PortBindings[appPort], network.PortBinding{
			HostIP:   netip.IPv4Unspecified(),
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
	cont, err := b.c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
		Platform:         nil,
		Name:             containerName(app.App.ID),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	_, err = b.c.ContainerStart(ctx, cont.ID, client.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start container")
	}
	return nil
}

func (b *Backend) synchronizeRuntime(ctx context.Context, apps []*domain.RuntimeDesiredState) error {
	// List old resources
	oldContainers, err := b.c.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: make(client.Filters).Add("label", fmt.Sprintf("%s=true", appLabel)),
	})
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}
	oldContainersMap := lo.SliceToMap(oldContainers.Items, func(c container.Summary) (string, *container.Summary) {
		return c.Labels[appIDLabel], &c
	})

	// Calculate next resources and apply
	for _, app := range apps {
		err = b.syncAppContainer(ctx, app, oldContainersMap[app.App.ID])
		if err != nil {
			slog.ErrorContext(ctx, "failed to sync app", "error", err)
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
	for _, oldContainer := range oldContainers.Items {
		appID := oldContainer.Labels[appIDLabel]
		if newApps[appID] {
			continue
		}

		_, err = b.c.ContainerRemove(ctx, oldContainer.ID, client.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to remove old container", "error", err)
			continue // fail-safe
		}
	}

	return nil
}
