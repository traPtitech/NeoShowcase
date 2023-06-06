package caddy

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/friendsofgo/errors"

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
}
`

func (s *server) Reconcile(sites []*domain.StaticSite) error {
	var b bytes.Buffer
	b.WriteString(":80 {\n")
	for _, site := range sites {
		matcherName := fmt.Sprintf("nsapp-%v", site.Application.ID)
		b.WriteString(fmt.Sprintf(
			siteTemplate,
			matcherName,
			web.HeaderNameSSGenAppID, site.Application.ID,
			matcherName,
			filepath.Join(s.c.DocsRoot, site.ArtifactID),
		))
	}
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
	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return errors.Errorf("expected 2xx, invalid status code %v", res.StatusCode)
	}
	return nil
}
