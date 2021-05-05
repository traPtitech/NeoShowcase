package main

import (
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
)

const (
	ModeDocker = iota
	ModeK8s
)

type Config struct {
	Mode    string                             `mapstructure:"mode" yaml:"mode"`
	Builder grpc.BuilderServiceClientConfig    `mapstructure:"builder" yaml:"builder"`
	SSGen   grpc.StaticSiteServiceClientConfig `mapstructure:"ssgen" yaml:"ssgen"`
	DB      admindb.Config                     `mapstructure:"db" yaml:"db"`
	HTTP    struct {
		Debug bool `mapstructure:"debug" yaml:"debug"`
		Port  int  `mapstructure:"port" yaml:"port"`
	} `mapstructure:"http" yaml:"http"`
	Image struct {
		Registry   string `mapstructure:"registry" yaml:"registry"`
		NamePrefix string `mapstructure:"namePrefix" yaml:"namePrefix"`
	} `mapstructure:"image" yaml:"image"`
}

func (c *Config) GetMode() int {
	switch strings.ToLower(c.Mode) {
	case "k8s", "kubernetes":
		return ModeK8s
	case "docker":
		return ModeDocker
	default:
		return ModeDocker
	}
}
