package main

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
)

type Config struct {
	ArtifactsRoot string `mapstructure:"artifactsRoot" yaml:"artifactsRoot"`
	BuiltIn       struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"builtIn" yaml:"builtIn"`
	NS      grpc.ComponentServiceClientConfig `mapstructure:"ns" yaml:"ns"`
	DB      admindb.Config                    `mapstructure:"db" yaml:"db"`
	Storage domain.StorageConfig              `mapstructure:"storage" yaml:"storage"`
}
