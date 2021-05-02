package builder

import (
	"github.com/traPtitech/neoshowcase/pkg/common"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
)

type Config struct {
	Buildkit struct {
		Address  string `mapstructure:"address" yaml:"address"`
		Registry string `mapstructure:"registry" yaml:"registry"`
	} `mapstructure:"buildkit" yaml:"buildkit"`
	GRPC    common.GRPCConfig    `mapstructure:"grpc" yaml:"grpc"`
	DB      admindb.Config       `mapstructure:"db" yaml:"db"`
	Storage common.StorageConfig `mapstructure:"storage" yaml:"storage"`
}
