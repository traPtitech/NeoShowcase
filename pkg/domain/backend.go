package domain

import (
	"context"
	"io"
	"strings"
)

type DesiredState struct {
	Runtime     []*RuntimeDesiredState
	StaticSites []*StaticSite
}

type DesiredStateLeader struct {
	TLSTargetDomains []string
}

type RuntimeDesiredState struct {
	App       *Application
	ImageName string
	ImageTag  string
	Envs      map[string]string
}

type ContainerEvent struct {
	ApplicationID string
}

type Container struct {
	ApplicationID string
	State         ContainerState
	Message       string
}

type ContainerState int

const (
	// ContainerStateMissing indicates that the container is not running.
	ContainerStateMissing ContainerState = iota
	// ContainerStateStarting indicates that the container is starting.
	ContainerStateStarting
	// ContainerStateRestarting indicates that the container is restarting.
	ContainerStateRestarting
	// ContainerStateRunning indicates that the container is running.
	ContainerStateRunning
	// ContainerStateExited indicates that the container has exited with code 0.
	ContainerStateExited
	// ContainerStateErrored indicates that the container has exited with a non-zero code.
	ContainerStateErrored
	// ContainerStateUnknown indicates that the container state is unknown.
	ContainerStateUnknown
)

type WildcardDomains []string

func (wd WildcardDomains) Validate() error {
	for _, d := range wd {
		if err := ValidateWildcardDomain(d); err != nil {
			return err
		}
	}
	return nil
}

func (wd WildcardDomains) TLSTargetDomain(website *Website) string {
	for _, d := range wd {
		if ContainsDomain(d, website.FQDN) {
			websiteParts := strings.Split(website.FQDN, ".")
			websiteParts[0] = "*"
			return strings.Join(websiteParts, ".")
		}
	}
	return website.FQDN
}

type Backend interface {
	Start(ctx context.Context) error
	Dispose(ctx context.Context) error

	AvailableDomains() AvailableDomainSlice
	TLSTargetDomain(website *Website) (host string, ok bool)
	AvailablePorts() AvailablePortSlice
	ListenContainerEvents() (sub <-chan *ContainerEvent, unsub func())
	Synchronize(ctx context.Context, s *DesiredState) error
	SynchronizeShared(ctx context.Context, s *DesiredStateLeader) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)
	AttachContainer(ctx context.Context, appID string, stdin io.Reader, stdout, stderr io.Writer) error
	ExecContainer(ctx context.Context, appID string, cmd []string, stdin io.Reader, stdout, stderr io.Writer) error
}
