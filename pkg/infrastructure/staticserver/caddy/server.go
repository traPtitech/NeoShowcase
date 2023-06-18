package caddy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/friendsofgo/errors"
	"github.com/samber/lo"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/web"
)

type Config struct {
	AdminAPI string `mapstructure:"adminAPI" yaml:"adminAPI"`
	DocsRoot string `mapstructure:"docsRoot" yaml:"docsRoot"`
}

type server struct {
	c Config
}

func NewServer(c Config) domain.StaticServer {
	return &server{c: c}
}

func (s *server) Start(_ context.Context) error {
	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	return nil
}

const siteTemplate = `@%v {
	header %v %v
}
file_server @%v {
	root %v
	precompressed br gzip zstd
}
`

const siteTemplateSPA = `@%v {
	header %v %v
	file {
		root %v
		try_files {path} {path}.html {path}/ {path}/index.html /index.html =404
	}
}
rewrite @%v {file_match.relative}
file_server @%v {
	root %v
	precompressed br gzip zstd
}
`

// unityWebglCompressionHeader for Unity WebGL https://docs.unity3d.com/Manual/webgl-deploying.html
const unityWebglCompressionHeader = `@gzip-suffix {
	path *.gz
}
@br-suffix {
	path *.br
}
header @gzip-suffix Content-Encoding gzip
header @br-suffix Content-Encoding br
`

func (s *server) Reconcile(sites []*domain.StaticSite) error {
	var b bytes.Buffer
	b.WriteString(":80 {\n")
	sites = lo.UniqBy(sites, func(site *domain.StaticSite) string { return site.Application.ID })
	for _, site := range sites {
		matcherName := fmt.Sprintf("nsapp-%v", site.Application.ID)
		if site.SPA {
			b.WriteString(fmt.Sprintf(
				siteTemplateSPA,
				matcherName,
				web.HeaderNameSSGenAppID, site.Application.ID,
				filepath.Join(s.c.DocsRoot, site.ArtifactID),
				matcherName,
				matcherName,
				filepath.Join(s.c.DocsRoot, site.ArtifactID),
			))
		} else {
			b.WriteString(fmt.Sprintf(
				siteTemplate,
				matcherName,
				web.HeaderNameSSGenAppID, site.Application.ID,
				matcherName,
				filepath.Join(s.c.DocsRoot, site.ArtifactID),
			))
		}
	}
	b.WriteString(unityWebglCompressionHeader)
	b.WriteString("}\n")
	return s.postConfig(b.Bytes())
}

func (s *server) postConfig(b []byte) error {
	req, err := http.NewRequest("POST", s.c.AdminAPI+"/load", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/caddyfile")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		resBody, _ := io.ReadAll(res.Body)
		return errors.Errorf("expected 2xx, invalid status code %v: %v", res.StatusCode, string(resBody))
	}
	return nil
}
