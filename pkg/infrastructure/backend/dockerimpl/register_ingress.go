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
	svcName := serviceName(website)
	netName := networkName(app.ID)

	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, traefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, traefikHTTPEntrypoint)
	}

	var middlewareNames []string
	switch app.Config.Authentication {
	case domain.AuthenticationTypeSoft:
		middlewareNames = append(middlewareNames,
			traefikAuthSoftMiddleware,
			traefikAuthMiddleware,
		)
	case domain.AuthenticationTypeHard:
		middlewareNames = append(middlewareNames,
			traefikAuthHardMiddleware,
			traefikAuthMiddleware,
		)
	}

	var middlewares m

	var rule string
	if website.PathPrefix == "/" {
		rule = fmt.Sprintf("Host(`%s`)", website.FQDN)
	} else {
		rule = fmt.Sprintf("Host(`%s`) && PathPrefix(`%s`)", website.FQDN, website.PathPrefix)
		middlewareName := stripMiddlewareName(website)
		middlewareNames = append(middlewareNames, middlewareName)
		middlewares[middlewareName] = m{
			"stripPrefix": m{
				"prefixes": []string{website.PathPrefix},
			},
		}
	}

	router := m{
		"entrypoints": entrypoints,
		"middlewares": middlewareNames,
		"service":     svcName,
		"rule":        rule,
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
			"middlewares": middlewares,
			"services": m{
				svcName: svc,
			},
		},
	}

	file, err := os.OpenFile(b.configFile(website), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}
