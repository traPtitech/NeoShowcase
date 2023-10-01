package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	configFilePath string
	config         Config
)

var rootCommand = &cobra.Command{
	Use:              "ns",
	Short:            "NeoShowcase",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cli.PrintVersion,
}

type component interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type componentGenF = func(c Config) (component, error)

func componentCommand(name string, gen componentGenF, longDesc string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("NeoShowcase %s component", name),
		Long:  longDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := gen(config)
			if err != nil {
				return err
			}

			go func() {
				err := service.Start(context.Background())
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("failed to start service: %+v", err)
				}
			}()
			log.Infof("NeoShowcase %s started", name)

			cli.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return service.Shutdown(ctx)
		},
	}
}

func main() {
	cobra.OnInitialize(cli.CobraOnInitializeFunc(&configFilePath, &config))

	rootCommand.AddCommand(
		componentCommand("auth-dev", NewAuthDev, ""),
		componentCommand("builder", NewBuilder, ""),
		componentCommand("controller", NewController, ""),
		componentCommand("gateway", NewGateway, ""),
		componentCommand("gitea-integration", NewGiteaIntegration, `Synchronizes gitea user / organization repositories and its owners with configured interval.
Operator needs to ensure that usernames of NeoShowcase and gitea are equivalent via SSO or some other method.
Admin token required.`),
		componentCommand("ssgen", NewSSGen, ""),
		cli.PrintConfCommand(&config),
	)

	flags := rootCommand.PersistentFlags()
	flags.StringVarP(&configFilePath, "config", "c", "", "config file path")
	cli.SetupDebugFlag(flags)
	cli.SetupLogLevelFlag(flags)

	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to exec: %+v", err)
	}
}
