package backend

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type Backend interface {
	CreateContainer(ctx context.Context, args domain.ContainerCreateArgs) error
	RestartContainer(ctx context.Context, appID string, envID string) error
	DestroyContainer(ctx context.Context, appID string, envID string) error
	ListContainers(ctx context.Context) ([]domain.Container, error)
	Dispose(ctx context.Context) error
}
