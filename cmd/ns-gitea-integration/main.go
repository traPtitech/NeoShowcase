package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	configFilePath string
	c              Config
)

var rootCommand = &cobra.Command{
	Use:   "ns-gitea integration",
	Short: "NeoShowcase Gitea Integration (optional component)",
	Long: `Synchronizes gitea user / organization repositories and its owners with configured interval.
Operator needs to ensure that usernames of NeoShowcase and gitea are equivalent via SSO or some other method.
Admin token required.`,
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cli.PrintVersion,
}

func runCommand() *cobra.Command {
	return &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := NewServer(c)
			if err != nil {
				return err
			}

			if c.Debug {
				boil.DebugMode = true
			}

			go func() {
				err := service.Start(context.Background())
				if err != nil {
					log.Fatalf("failed to start service: %+v", err)
				}
			}()

			log.Info("NeoShowcase Gitea Integration started")
			cli.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return service.Shutdown(ctx)
		},
	}
}

func main() {
	cobra.OnInitialize(cli.CobraOnInitializeFunc(&configFilePath, &c))

	rootCommand.AddCommand(
		runCommand(),
		cli.PrintConfCommand(&c),
	)

	flags := rootCommand.PersistentFlags()
	flags.StringVarP(&configFilePath, "config", "c", "", "config file path")
	cli.SetupDebugFlag(flags)
	cli.SetupLogLevelFlag(flags)

	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to exec: %+v", err)
	}
}
