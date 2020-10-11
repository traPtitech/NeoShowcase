package builder

import "github.com/traPtitech/neoshowcase/pkg/common"

type Config struct {
	Buildkit struct {
		Address string `mapstructure:"address" yaml:"address"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	GRPC common.GRPCConfig `mapstructure:"grpc" yaml:"grpc"`
	DB   common.DBConfig   `mapstructure:"db" yaml:"db"`
}
