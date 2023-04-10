package domain

import (
	"context"
)

type AppDesiredState struct {
	App       *Application
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
	ContainerStateMissing ContainerState = iota
	ContainerStateStarting
	ContainerStateRunning
	ContainerStateExited
	ContainerStateErrored
	ContainerStateUnknown
)

func TLSTargetDomain(allowWildcard bool, website *Website, ads AvailableDomainSlice) string {
	if allowWildcard {
		ad := ads.GetAvailableMatch(website.FQDN)
		if ad != nil {
			return ad.Domain
		}
	}
	return website.FQDN
}

type Backend interface {
	Start(ctx context.Context) error
	Dispose(ctx context.Context) error

	SynchronizeRuntime(ctx context.Context, apps []*AppDesiredState, ads AvailableDomainSlice) error
	SynchronizeSSIngress(ctx context.Context, sites []*StaticSite, ads AvailableDomainSlice) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)
}
