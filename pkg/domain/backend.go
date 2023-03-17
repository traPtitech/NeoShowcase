package domain

import (
	"context"
	"errors"
)

var (
	ErrContainerNotFound = errors.New("container not found")
)

type ContainerCreateArgs struct {
	ImageName string
	ImageTag  string
	Labels    map[string]string
	Envs      map[string]string
	Recreate  bool
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
	CreateContainer(ctx context.Context, app *Application, args ContainerCreateArgs) error
	RestartContainer(ctx context.Context, appID string) error
	DestroyContainer(ctx context.Context, app *Application) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]Container, error)
	Dispose(ctx context.Context) error
}
