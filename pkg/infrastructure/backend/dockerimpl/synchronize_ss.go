package dockerimpl

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type ssConfigBuilder struct {
	routers     m
	middlewares m
}

func newSSConfigBuilder() *ssConfigBuilder {
	return &ssConfigBuilder{
		routers:     make(m),
		middlewares: make(m),
	}
}

func (b *ssConfigBuilder) addStaticSite(site *domain.StaticSite) {
	router, newMiddlewares := routerBase(site.Application, site.Website, traefikSSServiceName)
	for name, mw := range newMiddlewares {
		b.middlewares[name] = mw
	}

	middlewareName := ssHeaderMiddlewareName(site)
	router["middlewares"] = append(router["middlewares"].([]string), middlewareName)
	b.middlewares[middlewareName] = m{
		"headers": m{
			"customRequestHeaders": m{
				web.HeaderNameSSGenAppID: site.Application.ID,
			},
		},
	}

	b.routers[traefikName(site.Website)] = router
}

func (b *ssConfigBuilder) build(ssURL string) m {
	http := m{
		"services": m{
			traefikSSServiceName: m{
				"loadBalancer": m{
					"servers": a{
						m{"url": ssURL},
					},
				},
			},
		},
	}
	if len(b.routers) > 0 {
		http["routers"] = b.routers
	}
	if len(b.middlewares) > 0 {
		http["middlewares"] = b.middlewares
	}
	return m{
		"http": http,
	}
}

func (b *dockerBackend) SynchronizeSSIngress(_ context.Context, sites []*domain.StaticSite) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	cb := newSSConfigBuilder()
	for _, site := range sites {
		cb.addStaticSite(site)
	}

	return b.writeConfig(traefikSSFilename, cb.build(b.ssURL))
}
