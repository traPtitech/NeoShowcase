package main

import (
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"google.golang.org/grpc"

	"strings"
)

const (
	ModeDocker = iota
	ModeK8s
)

type Config struct {
	Mode    string           `mapstructure:"mode" yaml:"mode"`
	Builder GRPCClientConfig `mapstructure:"builder" yaml:"builder"`
	SSGen   GRPCClientConfig `mapstructure:"ssgen" yaml:"ssgen"`
	DB      admindb.Config   `mapstructure:"db" yaml:"db"`
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

type GRPCClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func (c *GRPCClientConfig) Connect() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.Insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(c.Addr, opts...)
}
