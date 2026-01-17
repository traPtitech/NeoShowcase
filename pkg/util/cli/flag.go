package cli

import (
	"log/slog"

	"github.com/aarondl/sqlboiler/v4/boil"
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
			slog.SetLogLoggerLevel(slog.LevelDebug)
			boil.DebugMode = true
		}
	})
}

func SetupLogLevelFlag(flags *pflag.FlagSet) {
	flags.String("loglevel", "info", "log level (debug, info, warn, error)")
	BindPFlag(flags, "loglevel")
	viper.SetDefault("loglevel", "info")
	cobra.OnInitialize(func() {
		levelStr := viper.GetString("loglevel")
		var level slog.Level
		switch levelStr {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			slog.Error("invalid log level", "level", levelStr)
			level = slog.LevelInfo
		}
		slog.SetLogLoggerLevel(level)
	})
}
