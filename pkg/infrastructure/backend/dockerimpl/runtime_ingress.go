package dockerimpl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type (
	m map[string]any
	a []any
)

func routerBase(app *domain.Application, website *domain.Website, middlewares m) m {
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
		"rule":        rule,
	}

	if website.HTTPS {
		router["tls"] = m{
			"certResolver": traefikCertResolver,
			"domains": a{
				m{"main": website.FQDN},
			},
		}
	}
	return router
}

func (b *dockerBackend) runtimeIngressConfig(app *domain.Application, website *domain.Website) m {
	middlewares := m{}
	router := routerBase(app, website, middlewares)
	svcName := traefikName(website)
	router["service"] = svcName

	netName := networkName(app.ID)
	svc := m{
		"loadBalancer": m{
			"servers": a{
				m{"url": fmt.Sprintf("http://%s:%d/", netName, website.HTTPPort)},
			},
		},
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

	if len(middlewares) > 0 {
		config["http"].(m)["middlewares"] = middlewares
	}

	return config
}

func (b *dockerBackend) writeConfig(filename string, config any) error {
	file, err := os.OpenFile(filepath.Join(b.ingressConfDir, filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(config)
}

func (b *dockerBackend) synchronizeRuntimeIngresses(_ context.Context, app *domain.Application) error {
	entries, err := os.ReadDir(b.ingressConfDir)
	if err != nil {
		return err
	}

	confFilePrefix := configFilePrefix(app)
	existingFiles := m{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), confFilePrefix) {
			existingFiles[entry.Name()] = true
		}
	}

	newConfigs := m{}
	for _, website := range app.Websites {
		newConfigs[configFile(app, website)] = b.runtimeIngressConfig(app, website)
	}

	// Create / update configuration files
	for filename, config := range newConfigs {
		err = b.writeConfig(filename, config)
		if err != nil {
			return err
		}
	}

	// Prune old configuration files
	danglingFiles := util.MapDiff(existingFiles, newConfigs)
	for filename := range danglingFiles {
		err = os.Remove(filepath.Join(b.ingressConfDir, filename))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *dockerBackend) destroyRuntimeIngresses(_ context.Context, app *domain.Application) error {
	entries, err := os.ReadDir(b.ingressConfDir)
	if err != nil {
		return fmt.Errorf("failed to read conf dir: %w", err)
	}

	confFilePrefix := configFilePrefix(app)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), confFilePrefix) {
			err = os.Remove(filepath.Join(b.ingressConfDir, entry.Name()))
			if err != nil {
				return fmt.Errorf("failed to remove config: %w", err)
			}
		}
	}

	return nil
}
