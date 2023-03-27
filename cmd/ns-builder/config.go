package main

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
)

type Config struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	NS      grpc.ComponentServiceClientConfig `mapstructure:"ns" yaml:"ns"`
	DB      admindb.Config                    `mapstructure:"db" yaml:"db"`
	Storage domain.StorageConfig              `mapstructure:"storage" yaml:"storage"`
}
