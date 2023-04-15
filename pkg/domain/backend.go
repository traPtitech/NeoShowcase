package domain

import (
	"context"
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
		if MatchDomain(d, website.FQDN) {
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

	Synchronize(ctx context.Context, s *DesiredState) error
	GetContainer(ctx context.Context, appID string) (*Container, error)
	ListContainers(ctx context.Context) ([]*Container, error)
}
