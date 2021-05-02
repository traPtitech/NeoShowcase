package staticsitegen

import (
	"fmt"
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/common"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/webserver"
)

type Config struct {
	Server           string `mapstructure:"server" yaml:"server"`
	ArtifactsRoot    string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	GeneratedConfDir string `mapstructure:"generatedConfDir" yaml:"generatedConfDir"`
	Caddy            struct {
		AdminEndpoint string `mapstructure:"adminEndpoint" yaml:"adminEndpoint"`
	} `mapstructure:"caddy" yaml:"caddy"`
	BuiltIn struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"builtIn" yaml:"builtIn"`
	GRPC    common.GRPCConfig    `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config       `mapstructure:"db" yaml:"db"`
	Storage common.StorageConfig `mapstructure:"storage" yaml:"storage"`
}

func (c *Config) GetEngine() (webserver.Engine, error) {
	switch strings.ToLower(c.Server) {
	case "builtin":
		return &webserver.BuiltIn{
			ArtifactsRootPath: c.ArtifactsRoot,
			Port:              c.BuiltIn.Port,
		}, nil
	case "caddy":
		return &webserver.Caddy{
			ArtifactsRootPath: c.ArtifactsRoot,
			AdminEndpoint:     c.Caddy.AdminEndpoint,
		}, nil
	default:
		return nil, fmt.Errorf("unknown server type: %s", c.Server)
	}
}
