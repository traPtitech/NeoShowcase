package dockerimpl

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"maps"
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

func (b *ssConfigBuilder) addStaticSite(backend *Backend, site *domain.StaticSite) {
	router, newMiddlewares := backend.routerBase(site.Application, site.Website, traefikSSServiceName)
	maps.Copy(b.middlewares, newMiddlewares)

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

func (b *Backend) synchronizeSSIngress(_ context.Context, sites []*domain.StaticSite) error {
	cb := newSSConfigBuilder()
	for _, site := range sites {
		cb.addStaticSite(b, site)
	}
	return b.writeConfig(traefikSSFilename, cb.build(b.config.SS.URL))
}
