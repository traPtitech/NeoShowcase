package dockerimpl

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *dockerBackend) CreateContainer(ctx context.Context, app *domain.Application, args domain.ContainerCreateArgs) error {
	if args.ImageTag == "" {
		args.ImageTag = "latest"
	}

	// ビルドしたイメージをリポジトリからPull
	if err := b.c.PullImage(docker.PullImageOptions{
		Repository: args.ImageName,
		Tag:        args.ImageTag,
		Context:    ctx,
	}, docker.AuthConfiguration{}); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	var envs []string

	for name, value := range args.Envs {
		envs = append(envs, name+"="+value)
	}

	if args.Recreate {
		// 前のものが起動中の場合は削除する
		err := b.c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            containerName(app.ID),
			RemoveVolumes: true,
			Force:         true,
			Context:       ctx,
		})
		if err != nil {
			if _, ok := err.(*docker.NoSuchContainer); !ok {
				return fmt.Errorf("failed to remove old container: %w", err)
			}
		}
	}

	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: app.ID,
	})

	labels["traefik.enable"] = "true"
	for _, website := range app.Websites {
		traefikName := "nsapp_" + strings.ReplaceAll(website.FQDN, ".", "_")
		labels[fmt.Sprintf("traefik.http.routers.%s.rule", traefikName)] = fmt.Sprintf("Host(`%s`)", website.FQDN)
		labels[fmt.Sprintf("traefik.http.routers.%s.service", traefikName)] = traefikName
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", traefikName)] = strconv.Itoa(website.HTTPPort)
		switch app.Config.Authentication {
		case domain.AuthenticationTypeSoft:
			labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", traefikName)] = "ns_auth_soft@file, ns_auth@file"
		case domain.AuthenticationTypeHard:
			labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", traefikName)] = "ns_auth_hard@file, ns_auth@file"
		}
	}

	// ビルドしたイメージのコンテナを作成
	cont, err := b.c.CreateContainer(docker.CreateContainerOptions{
		Name: containerName(app.ID),
		Config: &docker.Config{
			Image:  args.ImageName + ":" + args.ImageTag,
			Labels: labels,
			Env:    envs,
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.RestartOnFailure(5),
		},
		NetworkingConfig: &docker.NetworkingConfig{EndpointsConfig: map[string]*docker.EndpointConfig{appNetwork: {}}},
		Context:          ctx,
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// コンテナを起動
	if err := b.c.StartContainer(cont.ID, nil); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}
