package builtin

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type Config struct {
	Port int `mapstructure:"port" yaml:"port"`
}

type server struct {
	port   int
	server *echo.Echo
	sites  map[string]*host

	sitesLock     sync.RWMutex
	reconcileLock sync.Mutex
}

type host struct {
	Echo *echo.Echo
	Site *domain.StaticSite
}

func NewServer(c Config) domain.StaticServer {
	b := &server{
		port:  c.Port,
		sites: make(map[string]*host),
	}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Any("/*", func(c echo.Context) error {
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
	b.server = e

	return b
}

func (b *server) Start(_ context.Context) error {
	return b.server.Start(fmt.Sprintf(":%d", b.port))
}

func (b *server) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

func (b *server) Reconcile(docsRoot string, sites []*domain.StaticSite) error {
	b.reconcileLock.Lock()
	defer b.reconcileLock.Unlock()

	siteMap := map[string]*host{}
	for _, site := range sites {
		artifactDir := filepath.Join(docsRoot, site.ArtifactID)
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
