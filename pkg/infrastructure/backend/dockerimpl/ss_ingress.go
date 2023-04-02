package dockerimpl

import (
	"context"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"gopkg.in/yaml.v2"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

func (b *dockerBackend) ReloadSSIngress(_ context.Context) error {
	b.reloadLock.Lock()
	defer b.reloadLock.Unlock()

	sites, err := domain.GetActiveStaticSites(context.Background(), b.appRepo, b.buildRepo)
	if err != nil {
		return err
	}

	routers := m{}
	middlewares := m{}

	for _, ss := range sites {
		router := routerBase(ss.Application, ss.Website, middlewares)
		router["service"] = traefikSSServiceName

		middlewareName := ssHeaderMiddlewareName(ss)
		router["middlewares"] = append(router["middlewares"].([]string), middlewareName)
		middlewares[middlewareName] = m{
			"headers": m{
				"customRequestHeaders": m{
					web.HeaderNameSSGenAppID: ss.Application.ID,
				},
			},
		}

		routers[traefikName(ss.Website)] = router
	}

	http := m{
		"services": m{
			traefikSSServiceName: m{
				"loadBalancer": m{
					"servers": a{
						m{"url": b.ssURL},
					},
				},
			},
		},
	}
	if len(routers) > 0 {
		http["routers"] = routers
	}
	if len(middlewares) > 0 {
		http["middlewares"] = middlewares
	}
	config := m{
		"http": http,
	}

	file, err := os.OpenFile(filepath.Join(b.ingressConfDir, traefikSSFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}
