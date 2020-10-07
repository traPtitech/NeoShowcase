package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var rootCommand = &cobra.Command{
	Use:     "ns-ssgen",
	Short:   "NeoShowcase StaticSiteGenerator",
	Version: fmt.Sprintf("%s (%s)", version, revision),
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func init() {
	rootCommand.AddCommand(
		runCommand(),
	)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
