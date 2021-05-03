package main

import (
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Config struct {
	ArtifactsRoot string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	BuiltIn       struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"builtIn" yaml:"builtIn"`
	GRPC struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config `mapstructure:"db" yaml:"db"`
	Storage storage.Config `mapstructure:"storage" yaml:"storage"`
}
