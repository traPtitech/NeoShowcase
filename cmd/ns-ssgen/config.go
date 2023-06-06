package main

import (
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/grpc"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/builtin"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/caddy"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
)

type Config struct {
	ArtifactsRoot string           `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	HealthPort    healthcheck.Port `mapstructure:"healthPort" yaml:"healthPort"`
	Server        struct {
		Type    string         `mapstructure:"type" yaml:"type"`
		BuiltIn builtin.Config `mapstructure:"builtIn" yaml:"builtIn"`
		Caddy   caddy.Config   `mapstructure:"caddy" yaml:"caddy"`
	} `mapstructure:"server" yaml:"server"`
	Controller grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	DB         repository.Config                  `mapstructure:"db" yaml:"db"`
	Storage    domain.StorageConfig               `mapstructure:"storage" yaml:"storage"`
}

func init() {
	viper.SetDefault("artifactsRoot", "/srv/artifacts")
	viper.SetDefault("healthPort", 8081)

	viper.SetDefault("server.type", "builtIn")
	viper.SetDefault("server.builtIn.port", 8080)
	viper.SetDefault("server.caddy.adminAPI", "http://localhost:2019")
	viper.SetDefault("server.caddy.configRoot", "/caddy-config")

	viper.SetDefault("controller.url", "http://ns-controller:10000")

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)

	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.dir", "/neoshowcase")
	viper.SetDefault("storage.s3.bucket", "neoshowcase")
	viper.SetDefault("storage.s3.endpoint", "")
	viper.SetDefault("storage.s3.region", "")
	viper.SetDefault("storage.s3.accessKey", "")
	viper.SetDefault("storage.s3.accessSecret", "")
	viper.SetDefault("storage.swift.username", "")
	viper.SetDefault("storage.swift.apiKey", "")
	viper.SetDefault("storage.swift.tenantName", "")
	viper.SetDefault("storage.swift.tenantId", "")
	viper.SetDefault("storage.swift.container", "neoshowcase")
	viper.SetDefault("storage.swift.authUrl", "")
}
