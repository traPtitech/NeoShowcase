package domain

import (
	"context"

	"github.com/volatiletech/null/v8"
)

type ContainerCreateArgs struct {
	ApplicationID string
	BranchID      string
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
	BranchID      string
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
	RestartContainer(ctx context.Context, appID string, branchID string) error
	DestroyContainer(ctx context.Context, appID string, branchID string) error
	ListContainers(ctx context.Context) ([]Container, error)
	RegisterIngress(ctx context.Context, appID string, branchID string, host string, destination null.String, port null.Int) error
	UnregisterIngress(ctx context.Context, appID string, branchID string) error
	Dispose(ctx context.Context) error
}
