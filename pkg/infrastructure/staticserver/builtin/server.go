package builtin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type Config struct {
	Port int `mapstructure:"port" yaml:"port"`
}

type server struct {
	docsRoot string
	server   *http.Server
	sites    map[string]*host

	sitesLock sync.RWMutex
}

type host struct {
	Echo *echo.Echo
	Site *domain.StaticSite
}

func NewServer(c Config, docsRoot string) domain.StaticServer {
	b := &server{
		docsRoot: docsRoot,
		sites:    make(map[string]*host),
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Any("/*", func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()

		appID := req.Header.Get(web.HeaderNameSSGenAppID)

		b.sitesLock.RLock()
		host := b.sites[appID]
		b.sitesLock.RUnlock()

		if host == nil {
			return echo.ErrNotFound
		}

		host.Echo.ServeHTTP(res, req)
		return nil
	})
	b.server = &http.Server{
		Addr:        fmt.Sprintf(":%d", int(c.Port)),
		Handler:     e,
		ReadTimeout: 30 * time.Second,
	}

	return b
}

func (b *server) Start(_ context.Context) error {
	if err := b.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (b *server) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

func (b *server) Reconcile(sites []*domain.StaticSite) error {
	siteMap := map[string]*host{}
	for _, site := range sites {
		if site.SPA {
			slog.Warn("SPA option is not supported in built-in static server")
		}
		artifactDir := filepath.Join(b.docsRoot, site.ArtifactID)
		e := echo.New()
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  artifactDir,
			Index: "index.html",
		}))
		siteMap[site.Application.ID] = &host{
			Echo: e,
			Site: site,
		}
	}

	b.sitesLock.Lock()
	b.sites = siteMap
	b.sitesLock.Unlock()
	return nil
}
