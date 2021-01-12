package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/traPtitech/neoshowcase/pkg/container"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (m *Manager) Create(ctx context.Context, args container.CreateArgs) (*container.CreateResult, error) {
	// ビルドしたイメージをリポジトリからPull
	if err := m.c.PullImage(docker.PullImageOptions{
		Repository: args.ImageName,
		Tag:        args.ImageTag,
		Context:    ctx,
	}, docker.AuthConfiguration{}); err != nil {
		return nil, fmt.Errorf("failed to pull image: %w", err)
	}

	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: args.ApplicationID,
	})

	if args.HTTPProxy != nil {
		labels = util.MergeLabels(labels, map[string]string{
			"traefik.enable": "true",
			fmt.Sprintf("traefik.http.routers.nsapp-%s.rule", args.ApplicationID):                      fmt.Sprintf("Host(`%s`)", args.HTTPProxy.Domain),
			fmt.Sprintf("traefik.http.services.nsapp-%s.loadbalancer.server.port", args.ApplicationID): fmt.Sprintf("%d", args.HTTPProxy.Port),
		})
	}

	// ビルドしたイメージのコンテナを作成
	cont, err := m.c.CreateContainer(docker.CreateContainerOptions{
		Name: containerName(args.ApplicationID),
		Config: &docker.Config{
			Image:  args.ImageName + ":" + args.ImageTag,
			Labels: labels,
			Env:    args.Envs,
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.RestartOnFailure(5),
		},
		NetworkingConfig: &docker.NetworkingConfig{EndpointsConfig: map[string]*docker.EndpointConfig{appNetwork: {}}},
		Context:          ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	if !args.NoStart {
		// コンテナを起動
		if err := m.c.StartContainer(cont.ID, nil); err != nil {
			return nil, fmt.Errorf("failed to start container: %w", err)
		}
	}
	return &container.CreateResult{}, nil
}
