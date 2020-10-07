package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	configFilePath string
)

var rootCommand = &cobra.Command{
	Use:     "ns",
	Short:   "NeoShowcase Core API Server",
	Version: fmt.Sprintf("%s (%s)", version, revision),
}

func runCommand() *cobra.Command {
	cmd := cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return &cmd
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetConfigFile(configFilePath)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetEnvPrefix("NS")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("failed to read config file: %v", err)
			}
		}
	})

	rootCommand.AddCommand(
		runCommand(),
	)

	flags := rootCommand.PersistentFlags()
	flags.StringVarP(&configFilePath, "config", "c", "./config.yaml", "config file path")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
