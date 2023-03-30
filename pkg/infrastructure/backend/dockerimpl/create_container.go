package dockerimpl

import (
	"context"

	"github.com/friendsofgo/errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
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
		return errors.Wrap(err, "failed to pull image")
	}

	envs := lo.MapToSlice(args.Envs, func(key string, value string) string {
		return key + "=" + value
	})

	// 前のものが起動中の場合は削除する
	err := b.c.RemoveContainer(docker.RemoveContainerOptions{
		ID:            containerName(app.ID),
		RemoveVolumes: true,
		Force:         true,
		Context:       ctx,
	})
	if err != nil {
		if _, ok := err.(*docker.NoSuchContainer); !ok {
			return errors.Wrap(err, "failed to remove old container")
		}
	}

	err = b.synchronizeRuntimeIngresses(ctx, app)
	if err != nil {
		return errors.Wrap(err, "failed to synchronize ingresses")
	}

	// ビルドしたイメージのコンテナを作成
	cont, err := b.c.CreateContainer(docker.CreateContainerOptions{
		Name: containerName(app.ID),
		Config: &docker.Config{
			Image:  args.ImageName + ":" + args.ImageTag,
			Labels: containerLabels(app.ID),
			Env:    envs,
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.RestartOnFailure(5),
		},
		NetworkingConfig: &docker.NetworkingConfig{EndpointsConfig: map[string]*docker.EndpointConfig{
			appNetwork: {
				Aliases: []string{networkName(app.ID)},
			},
		}},
		Context: ctx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create container")
	}

	// コンテナを起動
	if err := b.c.StartContainer(cont.ID, nil); err != nil {
		return errors.Wrap(err, "failed to start container")
	}

	return nil
}
