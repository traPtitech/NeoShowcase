package domain

import (
	"context"
	"errors"

	"github.com/volatiletech/null/v8"
)

var (
	ErrContainerNotFound = errors.New("container not found")
)

type ContainerCreateArgs struct {
	ApplicationID string
	ImageName     string
	ImageTag      string
	Labels        map[string]string
	Envs          map[string]string
	HTTPProxy     *ContainerHTTPProxy
	Recreate      bool
}

type ContainerHTTPProxy struct {
	Domain string
	Port   int
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
	CreateContainer(ctx context.Context, args ContainerCreateArgs) error
	RestartContainer(ctx context.Context, appID string) error
	DestroyContainer(ctx context.Context, appID string) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]Container, error)
	RegisterIngress(ctx context.Context, appID string, host string, destination null.String, port null.Int) error
	UnregisterIngress(ctx context.Context, appID string) error
	Dispose(ctx context.Context) error
}
