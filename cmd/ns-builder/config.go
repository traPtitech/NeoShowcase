package main

import (
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
)

type Config struct {
	Buildkit struct {
		Address  string `mapstructure:"address" yaml:"address"`
		Registry string `mapstructure:"registry" yaml:"registry"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	GRPC struct {
		Port int `mapstructure:"port" yaml:"port"`
	} `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config `mapstructure:"db" yaml:"db"`
	Storage storage.Config `mapstructure:"storage" yaml:"storage"`
}
