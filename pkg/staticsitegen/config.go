package staticsitegen

import (
	"fmt"
	"github.com/traPtitech/neoshowcase/pkg/common"
	"github.com/traPtitech/neoshowcase/pkg/staticsitegen/generator"
	"strings"
)

type Config struct {
	Server      string            `mapstructure:"server" yaml:"server"`
	NginxConfig generator.Nginx   `mapstructure:"nginx" yaml:"nginx"`
	GRPC        common.GRPCConfig `mapstructure:"grpc" yaml:"grpc"`
	DB          common.DBConfig   `mapstructure:"db" yaml:"db"`
}

func (c *Config) GetEngine() (generator.Engine, error) {
	switch strings.ToLower(c.Server) {
	case "nginx":
		return &c.NginxConfig, nil
	default:
		return nil, fmt.Errorf("unknown server type: %s", c.Server)
	}
}
