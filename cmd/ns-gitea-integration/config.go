package main

import (
	"github.com/spf13/viper"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/repository"
	giteaintegration "github.com/traPtitech/neoshowcase/pkg/usecase/gitea-integration"
)

type Config struct {
	Debug bool                    `mapstructure:"debug" yaml:"debug"`
	Gitea giteaintegration.Config `mapstructure:"gitea" yaml:"gitea"`
	DB    repository.Config       `mapstructure:"db" yaml:"db"`
}

func init() {
	viper.SetDefault("debug", false)

	viper.SetDefault("gitea.url", "https://git.trap.jp")
	viper.SetDefault("gitea.token", "")
	viper.SetDefault("gitea.intervalSeconds", 600)

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)
}
