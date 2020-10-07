package staticsitegen

import (
	"context"
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/apiserver/grpc/api"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
	"strings"
)

type Service struct {
	engine generator.Engine
	config Config
	client api.SitesClient
}

func New(c Config, client api.SitesClient) (*Service, error) {
	s := &Service{
		config: c,
		client: client,
	}

	// server type
	switch strings.ToLower(c.ServerType) {
	case "nginx":
		s.engine = &c.NginxConfig
	default:
		return nil, fmt.Errorf("unknown server type: %s", c.ServerType)
	}

	return s, nil
}

func (s *Service) Start(_ context.Context) error {
	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	return nil
}
