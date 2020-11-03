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
