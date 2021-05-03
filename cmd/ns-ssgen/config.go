package main

import (
	"github.com/traPtitech/neoshowcase/pkg/common"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
)

type Config struct {
	ArtifactsRoot string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	BuiltIn       struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"builtIn" yaml:"builtIn"`
	GRPC    common.GRPCConfig    `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config       `mapstructure:"db" yaml:"db"`
	Storage common.StorageConfig `mapstructure:"storage" yaml:"storage"`
}
