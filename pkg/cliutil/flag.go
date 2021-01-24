package cliutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Debug = false

func SetupDebugFlag(flags *pflag.FlagSet) {
	flags.BoolVar(&Debug, "debug", false, "debug mode")
	BindPFlag(flags, "debug")
	viper.SetDefault("debug", false)
	cobra.OnInitialize(func() {
		if Debug {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
		}
	})
}

func SetupLogLevelFlag(flags *pflag.FlagSet) {
	flags.BoolVar(&Debug, "loglevel", false, "log level (trace, debug, info, warn, error)")
	BindPFlag(flags, "loglevel")
	viper.SetDefault("loglevel", "info")
	cobra.OnInitialize(func() {
		level, err := log.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			log.Error(err.Error())
		} else {
			log.SetLevel(level)
		}
	})
}
