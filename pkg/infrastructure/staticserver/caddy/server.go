package caddy

import (
	"context"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

type Config struct {
	AdminAPI   string `mapstructure:"adminAPI" yaml:"adminAPI"`
	ConfigRoot string `mapstructure:"configRoot" yaml:"configRoot"`
}

type server struct {
	adminAPI   string
	configRoot string
}

func NewServer(c Config) domain.StaticServer {
	return &server{
		adminAPI:   c.AdminAPI,
		configRoot: c.ConfigRoot,
	}
}

func (s *server) Start(_ context.Context) error {
	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	return nil
}

func (s *server) Reconcile(docsRoot string, sites []*domain.StaticSite) error {
	panic("implement me")
}
