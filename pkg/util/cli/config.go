package cli

import (
	"strings"

	log "github.com/sirupsen/logrus"
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

func CobraOnInitializeFunc(configFilePath *string, config interface{}) func() {
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
				log.Fatalf("failed to read config file: %v", err)
			}
		}
		if err := viper.Unmarshal(config); err != nil {
			log.Fatalf("failed to unmarshal config: %v", err)
		}
	}
}
