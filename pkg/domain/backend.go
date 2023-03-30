package domain

import (
	"context"

	"github.com/friendsofgo/errors"
)

var (
	ErrContainerNotFound = errors.New("container not found")
)

type ContainerCreateArgs struct {
	ImageName string
	ImageTag  string
	Envs      map[string]string
}

type Container struct {
	ApplicationID string
	State         ContainerState
}

type ContainerState int

const (
	ContainerStateRunning ContainerState = iota
	ContainerStateRestarting
	ContainerStateStopped
	ContainerStateOther
)

type Backend interface {
	Start(ctx context.Context) error
	Dispose(ctx context.Context) error

	CreateContainer(ctx context.Context, app *Application, args ContainerCreateArgs) error
	DestroyContainer(ctx context.Context, app *Application) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]Container, error)
	ReloadSSIngress(ctx context.Context) error
}
