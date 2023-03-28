package main

import (
	"strings"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/backend/dockerimpl"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/dbmanager"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

const (
	ModeDocker = iota
	ModeK8s
)

type Config struct {
	Debug   bool                                  `mapstructure:"debug" yaml:"debug"`
	Mode    string                                `mapstructure:"mode" yaml:"mode"`
	SS      domain.StaticServerConnectivityConfig `mapstructure:"ss" yaml:"ss"`
	DB      admindb.Config                        `mapstructure:"db" yaml:"db"`
	MariaDB dbmanager.MariaDBConfig               `mapstructure:"mariadb" yaml:"mariadb"`
	MongoDB dbmanager.MongoDBConfig               `mapstructure:"mongodb" yaml:"mongodb"`
	Docker  struct {
		ConfDir string `mapstructure:"confDir" yaml:"confDir"`
	} `mapstructure:"docker" yaml:"docker"`
	GRPC struct {
		App struct {
			Port int `mapstructure:"port" yaml:"port"`
		} `mapstructure:"app" yaml:"app"`
		Component struct {
			Port int `mapstructure:"port" yaml:"port"`
		} `mapstructure:"component" yaml:"component"`
	} `mapstructure:"grpc" yaml:"grpc"`
	Repository struct {
		CacheDir string `mapstructure:"cacheDir" yaml:"cacheDir"`
	} `mapstructure:"repository" yaml:"repository"`
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

func provideIngressConfDirPath(c Config) dockerimpl.IngressConfDirPath {
	return dockerimpl.IngressConfDirPath(c.Docker.ConfDir)
}

func provideImageRegistry(c Config) builder.DockerImageRegistryString {
	return c.Image.Registry
}

func provideImagePrefix(c Config) builder.DockerImageNamePrefixString {
	return c.Image.NamePrefix
}

func provideRepositoryFetcherCacheDir(c Config) usecase.RepositoryFetcherCacheDir {
	return usecase.RepositoryFetcherCacheDir(c.Repository.CacheDir)
}
