package main

import (
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
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
	MariaDB dbmanager.MariaDBConfig            `mapstructure:"mariadb" yaml:"mariadb"`
	MongoDB dbmanager.MongoDBConfig            `mapstructure:"mongodb" yaml:"mongodb"`
	GRPC    struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"grpc" yaml:"grpc"`
	HTTP struct {
		Debug bool `mapstructure:"debug" yaml:"debug"`
		Port  int  `mapstructure:"port" yaml:"port"`
	} `mapstructure:"http" yaml:"http"`
	Image struct {
		Registry   builder.DockerImageRegistryString   `mapstructure:"registry" yaml:"registry"`
		NamePrefix builder.DockerImageNamePrefixString `mapstructure:"namePrefix" yaml:"namePrefix"`
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

func provideImageRegistry(c Config) builder.DockerImageRegistryString {
	return c.Image.Registry
}

func provideImagePrefix(c Config) builder.DockerImageNamePrefixString {
	return c.Image.NamePrefix
}
