package dockerimpl

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

func (b *dockerBackend) CreateContainer(ctx context.Context, args domain.ContainerCreateArgs) error {
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

	labels := util.MergeLabels(args.Labels, map[string]string{
		appContainerLabel:              "true",
		appContainerApplicationIDLabel: args.ApplicationID,
		appContainerEnvironmentIDLabel: args.EnvironmentID,
	})

	var envs []string

	for name, value := range args.Envs {
		envs = append(envs, name+"="+value)
	}

	if args.Recreate {
		// 前のものが起動中の場合は削除する
		err := b.c.RemoveContainer(docker.RemoveContainerOptions{
			ID:            containerName(args.ApplicationID, args.EnvironmentID),
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

	// ビルドしたイメージのコンテナを作成
	cont, err := b.c.CreateContainer(docker.CreateContainerOptions{
		Name: containerName(args.ApplicationID, args.EnvironmentID),
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
