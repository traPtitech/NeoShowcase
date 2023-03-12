package main

import (
	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
)

type Config struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	GRPC struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config       `mapstructure:"db" yaml:"db"`
	Storage domain.StorageConfig `mapstructure:"storage" yaml:"storage"`
}
