package backend

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/volatiletech/null/v8"
)

type Backend interface {
	CreateContainer(ctx context.Context, args domain.ContainerCreateArgs) error
	RestartContainer(ctx context.Context, appID string, envID string) error
	DestroyContainer(ctx context.Context, appID string, envID string) error
	ListContainers(ctx context.Context) ([]domain.Container, error)
	RegisterIngress(ctx context.Context, appID string, envID string, host string, destination null.String, port null.Int) error
	UnregisterIngress(ctx context.Context, appID string, envID string) error
	Dispose(ctx context.Context) error
}
