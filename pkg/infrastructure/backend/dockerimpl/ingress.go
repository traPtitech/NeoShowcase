package dockerimpl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type (
	m map[string]any
	a []any
)

func (b *Backend) routerBase(website *domain.Website, svcName string) (router m, middlewares m) {
	middlewares = make(m)

	var entrypoints []string
	if website.HTTPS {
		entrypoints = append(entrypoints, web.TraefikHTTPSEntrypoint)
	} else {
		entrypoints = append(entrypoints, web.TraefikHTTPEntrypoint)
	}

	var middlewareNames []string
	authConfig := b.targetAuth(website.FQDN)
	if authConfig != nil {
		switch website.Authentication {
		case domain.AuthenticationTypeSoft:
			middlewareNames = append(middlewareNames, authConfig.Soft...)
		case domain.AuthenticationTypeHard:
			middlewareNames = append(middlewareNames, authConfig.Hard...)
		}
	} else if website.Authentication != domain.AuthenticationTypeOff {
		log.Warnf("auth config not available for %s", website.FQDN)
	}

	var rule string
	if website.PathPrefix == "/" {
		rule = fmt.Sprintf("Host(`%s`)", website.FQDN)
	} else {
		rule = fmt.Sprintf("Host(`%s`) && PathPrefix(`%s`)", website.FQDN, website.PathPrefix)
		if website.StripPrefix {
			middlewareName := stripMiddlewareName(website)
			middlewareNames = append(middlewareNames, middlewareName)
			middlewares[middlewareName] = m{
				"stripPrefix": m{
					"prefixes": []string{website.PathPrefix},
				},
			}
		}
	}

	priorityOffset := b.config.Routing.Traefik.PriorityOffset
	priority := len(rule) + priorityOffset

	router = m{
		"entrypoints": entrypoints,
		"middlewares": middlewareNames,
		"rule":        rule,
		"priority":    priority,
		"service":     svcName,
	}

	if website.HTTPS {
		targetDomain := b.config.TLS.Wildcard.Domains.TLSTargetDomain(website)
		router["tls"] = m{
			"certResolver": b.config.TLS.CertResolver,
			"domains": a{
				m{"main": targetDomain},
			},
		}
	}

	return router, middlewares
}

type runtimeConfigBuilder struct {
	routers     m
	middlewares m
	services    m
}

func newRuntimeConfigBuilder() *runtimeConfigBuilder {
	return &runtimeConfigBuilder{
		routers:     make(m),
		middlewares: make(m),
		services:    make(m),
	}
}

func (b *runtimeConfigBuilder) addWebsite(backend *Backend, app *domain.Application, website *domain.Website) {
	svcName := traefikName(website)

	router, middlewares := backend.routerBase(website, svcName)

	netName := networkName(app.ID)
	svc := m{
		"loadBalancer": m{
			"servers": a{
				m{"url": fmt.Sprintf(
					"%s://%s:%d/",
					lo.Ternary(website.H2C, "h2c", "http"),
					netName,
					website.HTTPPort,
				)},
			},
		},
	}

	b.routers[svcName] = router
	for name, mw := range middlewares {
		b.middlewares[name] = mw
	}
	b.services[svcName] = svc
}

func (b *runtimeConfigBuilder) build() m {
	http := make(m)

	if len(b.routers) > 0 {
		http["routers"] = b.routers
	}
	if len(b.middlewares) > 0 {
		http["middlewares"] = b.middlewares
	}
	if len(b.services) > 0 {
		http["services"] = b.services
	}

	if len(http) == 0 {
		return nil
	}
	return m{
		"http": http,
	}
}

func (b *Backend) writeConfig(filename string, config any) error {
	file, err := os.OpenFile(filepath.Join(b.config.ConfDir, filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}
