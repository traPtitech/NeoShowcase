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
	ContainerStateMissing ContainerState = iota
	ContainerStateStarting
	ContainerStateRestarting
	ContainerStateRunning
	ContainerStateIdle
	ContainerStateExited
	ContainerStateErrored
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
	AvailablePorts() AvailablePortSlice
	ListenContainerEvents() (sub <-chan *ContainerEvent, unsub func())
	Synchronize(ctx context.Context, s *DesiredState) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)
	AttachContainer(ctx context.Context, appID string, stdin io.Reader, stdout, stderr io.Writer) error
	ExecContainer(ctx context.Context, appID string, cmd []string, stdin io.Reader, stdout, stderr io.Writer) error
}
