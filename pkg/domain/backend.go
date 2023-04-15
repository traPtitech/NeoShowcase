package domain

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/util/ds"
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

func (wd WildcardDomains) IsValid() bool {
	return lo.EveryBy(wd, IsValidWildcardDomain)
}

func (wd WildcardDomains) TLSTargetDomain(website *Website) string {
	websiteParts := strings.Split(website.FQDN, ".")
	for _, d := range wd {
		baseParts := strings.Split(strings.TrimPrefix(d, "*."), ".")
		if ds.HasSuffix(websiteParts, baseParts) {
			switch {
			case len(websiteParts) == len(baseParts)+1:
				return d
			case len(websiteParts) > len(baseParts)+1:
				websiteParts[0] = "*"
				return strings.Join(websiteParts, ".")
			}
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
