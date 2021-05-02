package webserver

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	storage2 "github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/util"
)

type BuiltIn struct {
	ArtifactsRootPath string
	Port              int

	storage   storage2.Storage
	server    *echo.Echo
	sites     map[string]*builtInHost
	sitesLock sync.RWMutex
}

type builtInHost struct {
	Echo *echo.Echo
	Site *Site
}

func (b *BuiltIn) Init(s storage2.Storage) error {
	b.storage = s
	b.sites = map[string]*builtInHost{}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Any("/*", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		b.sitesLock.RLock()
		host := b.sites[req.Host]
		b.sitesLock.RUnlock()

		if host == nil {
			return echo.ErrNotFound
		}

		host.Echo.ServeHTTP(res, req)
		return nil
	})
	b.server = e

	return nil
}

func (b *BuiltIn) Start(ctx context.Context) error {
	go func() {
		err := b.server.Start(fmt.Sprintf(":%d", b.Port))
		if err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()
	return nil
}

func (b *BuiltIn) Reconcile(sites []*Site) error {
	siteMap := map[string]*builtInHost{}
	for _, site := range sites {
		artifactDir := filepath.Join(b.ArtifactsRootPath, site.ArtifactID)

		e := echo.New()
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  artifactDir,
			Index: "index.html",
		}))
		siteMap[site.FQDN] = &builtInHost{
			Echo: e,
			Site: site,
		}

		// 静的ファイルの配置
		if !util.FileExists(artifactDir) {
			if err := storage2.ExtractTarToDir(b.storage, filepath.Join("artifacts", site.ArtifactID+".tar"), artifactDir); err != nil {
				return fmt.Errorf("failed to extract artifact tar: %w", err)
			}
		}
	}

	b.sitesLock.Lock()
	b.sites = siteMap
	b.sitesLock.Unlock()
	return nil
}

func (b *BuiltIn) Close(ctx context.Context) error {
	return b.server.Shutdown(ctx)
}
