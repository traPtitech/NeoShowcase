package dockerimpl

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type (
	m map[string]any
	a []any
)

func (b *dockerBackend) registerIngress(_ context.Context, app *domain.Application, website *domain.Website) error {
	svcName := serviceName(website.FQDN)
	netName := networkName(app.ID)

	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, traefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, traefikHTTPEntrypoint)
	}

	var middlewares []string
	switch app.Config.Authentication {
	case domain.AuthenticationTypeSoft:
		middlewares = append(middlewares,
			traefikAuthSoftMiddleware,
			traefikAuthMiddleware,
		)
	case domain.AuthenticationTypeHard:
		middlewares = append(middlewares,
			traefikAuthHardMiddleware,
			traefikAuthMiddleware,
		)
	}

	router := m{
		"entrypoints": entrypoints,
		"middlewares": middlewares,
		"service":     svcName,
		"rule":        fmt.Sprintf("Host(`%s`)", website.FQDN),
	}
	svc := m{
		"loadBalancer": m{
			"servers": a{
				m{"url": fmt.Sprintf("http://%s:%d/", netName, website.HTTPPort)},
			},
		},
	}

	if website.HTTPS {
		router["tls"] = m{
			"certResolver": traefikCertResolver,
			"domains": a{
				m{"main": website.FQDN},
			},
		}
	}

	config := m{
		"http": m{
			"routers": m{
				svcName: router,
			},
			"services": m{
				svcName: svc,
			},
		},
	}

	file, err := os.OpenFile(b.configFile(website.FQDN), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}
