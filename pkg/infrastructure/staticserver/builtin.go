package staticserver

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

type BuiltIn struct {
	docsRoot      string
	port          int
	server        *echo.Echo
	sites         map[string]*builtInHost
	sitesLock     sync.RWMutex
	reconcileLock sync.Mutex
}

type builtInHost struct {
	Echo *echo.Echo
	Site *domain.StaticSite
}

func NewBuiltIn(path domain.StaticServerDocumentRootPath, port domain.StaticServerPort) domain.SSEngine {
	b := &BuiltIn{
		docsRoot: string(path),
		port:     int(port),
		server:   nil,
		sites:    map[string]*builtInHost{},
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

func (b *BuiltIn) Start(_ context.Context) error {
	return b.server.Start(fmt.Sprintf(":%d", b.port))
}

func (b *BuiltIn) Shutdown(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}

func (b *BuiltIn) Reconcile(sites []*domain.StaticSite) error {
	b.reconcileLock.Lock()
	defer b.reconcileLock.Unlock()

	siteMap := map[string]*builtInHost{}
	for _, site := range sites {
		artifactDir := filepath.Join(b.docsRoot, site.ArtifactID)
		e := echo.New()
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  artifactDir,
			Index: "index.html",
		}))
		siteMap[site.Application.ID] = &builtInHost{
			Echo: e,
			Site: site,
		}
	}

	b.sitesLock.Lock()
	b.sites = siteMap
	b.sitesLock.Unlock()
	return nil
}
