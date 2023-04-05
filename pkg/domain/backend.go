package domain

import (
	"context"

	"github.com/friendsofgo/errors"
)

var (
	ErrContainerNotFound = errors.New("container not found")
)

type AppDesiredState struct {
	App       *Application
	ImageName string
	ImageTag  string
	Envs      map[string]string
	Restart   bool
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

	Synchronize(ctx context.Context, apps []*AppDesiredState) error
	SynchronizeSSIngress(ctx context.Context) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)
}
