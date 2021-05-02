package webserver

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	storage2 "github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Caddy struct {
	ArtifactsRootPath string
	AdminEndpoint     string
	tmpls             *template.Template
	storage           storage2.Storage
}

func (engine *Caddy) Init(s storage2.Storage) (err error) {
	engine.tmpls, err = template.New("caddyfile").Parse(`
{
  auto_https off
  admin 0.0.0.0:2019
}

{{ $rootDir := .ArtifactRoot }}

{{ range .Sites }}
http://{{ .FQDN }} {
	root * {{ $rootDir }}/{{ .ArtifactID }}
	file_server
}
{{ end }}
`)
	engine.storage = s
	return
}

func (engine *Caddy) Start(ctx context.Context) error {
	return nil
}

func (engine *Caddy) Reconcile(sites []*Site) error {
	var sitesData []map[string]interface{}
	for _, site := range sites {
		sitesData = append(sitesData, map[string]interface{}{
			"ArtifactID": site.ArtifactID,
			"FQDN":       site.FQDN,
		})

	}

	// 設定ファイル生成
	b := &bytes.Buffer{}
	if err := engine.tmpls.ExecuteTemplate(b, "caddyfile", map[string]interface{}{
		"Sites":        sitesData,
		"ArtifactRoot": strings.TrimRight(engine.ArtifactsRootPath, "/"),
	}); err != nil {
		return fmt.Errorf("failed to generate conf file: %w", err)
	}
	log.Debug(b.String())

	// caddy reload
	resp, err := http.Post(engine.AdminEndpoint+"/load", "text/caddyfile", b)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}
	defer resp.Body.Close()
	log.Debug(resp.StatusCode)

	return nil
}

func (engine *Caddy) Close(ctx context.Context) error {
	return nil
}
