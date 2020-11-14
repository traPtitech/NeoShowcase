package generator

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/traPtitech/neoshowcase/pkg/util"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

type Caddy struct {
	ArtifactsRootPath string
	AdminEndpoint     string
	tmpls             *template.Template
}

func (engine *Caddy) Init() (err error) {
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
	return
}

func (engine *Caddy) Reconcile(sites []*Site) error {
	var sitesData []map[string]interface{}
	for _, site := range sites {
		sitesData = append(sitesData, map[string]interface{}{
			"ArtifactID": site.ArtifactID,
			"FQDN":       site.FQDN,
		})

		// 静的ファイルの配置
		artifactDir := filepath.Join(engine.ArtifactsRootPath, site.ArtifactID)
		if !util.FileExists(artifactDir) {
			// TODO artifactのtarの取り出しをStorageインターフェース経由で抽象化
			if err := util.ExtractTarToDir(filepath.Join("/neoshowcase/artifacts", site.ArtifactID+".tar"), artifactDir); err != nil {
				return fmt.Errorf("failed to extract artifact tar: %w", err)
			}
		}
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
