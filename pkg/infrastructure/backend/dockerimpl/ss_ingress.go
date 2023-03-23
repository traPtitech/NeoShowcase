package dockerimpl

import (
	"context"
	"fmt"
	"os"

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

	config := m{
		"http": m{
			"routers":     routers,
			"middlewares": middlewares,
			"services": m{
				traefikSSServiceName: m{
					"loadBalancer": m{
						"servers": a{
							m{"url": b.ssURL},
						},
					},
				},
			},
		},
	}

	file, err := os.OpenFile(traefikSSFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}
