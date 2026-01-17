package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

func init() {
	cli.SetVersion(version, revision)
}

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

			ctx := cmd.Context()
			ctx, cancel := context.WithCancel(ctx)
			go func() {
				err := service.Start(ctx)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("failed to start service", "error", err)
					cancel()
				}
			}()
			slog.Info("NeoShowcase service started", "service", name)

			cli.WaitSIGINT()
			cancel()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return service.Shutdown(shutdownCtx)
		},
	}
}

func main() {
	cobra.OnInitialize(cli.CobraOnInitializeFunc(&configFilePath, &config))

	rootCommand.AddCommand(
		componentCommand("auth-dev", NewAuthDev, ""),
		componentCommand("builder", NewBuilder, ""),
		componentCommand("buildpack-helper", NewBuildpackHelper, ""),
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
		slog.Error("failed to exec", "error", err)
		os.Exit(1)
	}
}
