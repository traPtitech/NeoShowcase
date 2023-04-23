package main

import (
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
)

type Config struct {
	ArtifactsRoot string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	BuiltIn       struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"builtIn" yaml:"builtIn"`
	Controller grpc.ControllerServiceClientConfig `mapstructure:"controller" yaml:"controller"`
	DB         admindb.Config                     `mapstructure:"db" yaml:"db"`
	Storage    domain.StorageConfig               `mapstructure:"storage" yaml:"storage"`
}

func init() {
	viper.SetDefault("artifactsRoot", "/srv/artifacts")

	viper.SetDefault("builtIn.port", 8080)

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
