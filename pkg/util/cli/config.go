package cli

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindPFlag(flags *pflag.FlagSet, key string, flag ...string) {
	if len(flag) == 0 {
		flag = []string{key}
	}
	if err := viper.BindPFlag(key, flags.Lookup(flag[0])); err != nil {
		panic(err)
	}
}

func CobraOnInitializeFunc(configFilePath *string, config any) func() {
	return func() {
		if len(*configFilePath) > 0 {
			viper.SetConfigFile(*configFilePath)
		} else {
			viper.AddConfigPath(".")
			viper.SetConfigName("config")
		}
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetEnvPrefix("NS")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				slog.Error("failed to read config file", "error", err)
				os.Exit(1)
			}
		}
		if err := viper.Unmarshal(config); err != nil {
			slog.Error("failed to unmarshal config", "error", err)
			os.Exit(1)
		}
	}
}
