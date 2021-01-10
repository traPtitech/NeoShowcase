package apiserver

import (
	"github.com/traPtitech/neoshowcase/pkg/common"
	"strings"
)

const (
	ModeDocker = iota
	ModeK8s
)

type Config struct {
	Mode    string                  `mapstructure:"mode" yaml:"mode"`
	Builder common.GRPCClientConfig `mapstructure:"builder" yaml:"builder"`
	SSGen   common.GRPCClientConfig `mapstructure:"ssgen" yaml:"ssgen"`
	DB      common.DBConfig         `mapstructure:"db" yaml:"db"`
	HTTP    struct {
		Debug bool `mapstructure:"debug" yaml:"debug"`
		Port  int  `mapstructure:"port" yaml:"port"`
	} `mapstructure:"http" yaml:"http"`
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
