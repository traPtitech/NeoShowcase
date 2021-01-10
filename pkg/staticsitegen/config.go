package staticsitegen

import (
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/common"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
	"strings"
)

type Config struct {
	Server           string `mapstructure:"server" yaml:"server"`
	ArtifactsRoot    string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	GeneratedConfDir string `mapstructure:"generatedConfDir" yaml:"generatedConfDir"`
	Caddy            struct {
		AdminEndpoint string `mapstructure:"adminEndpoint" yaml:"adminEndpoint"`
	} `mapstructure:"caddy" yaml:"caddy"`
	GRPC    common.GRPCConfig    `mapstructure:"grpc" yaml:"grpc"`
	DB      common.DBConfig      `mapstructure:"db" yaml:"db"`
	Storage common.StorageConfig `mapstructure:"storage" yaml:"storage"`
}

func (c *Config) GetEngine() (generator.Engine, error) {
	switch strings.ToLower(c.Server) {
	case "caddy":
		return &generator.Caddy{
			ArtifactsRootPath: c.ArtifactsRoot,
			AdminEndpoint:     c.Caddy.AdminEndpoint,
		}, nil
	default:
		return nil, fmt.Errorf("unknown server type: %s", c.Server)
	}
}
