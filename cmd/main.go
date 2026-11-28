package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/traPtitech/neoshowcase/pkg/infrastructure/observability"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
	"github.com/traPtitech/neoshowcase/pkg/util/slogutil"
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
	Use:     "ns",
	Short:   "NeoShowcase",
	Version: fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		serviceName := cmd.Name()
		slogutil.InitLogger(serviceName, version, revision)
	},
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
			return runComponentUntilSignal(cmd.Context(), name, service)
		},
	}
}

// runComponentUntilSignal runs the service until the process receives SIGINT or SIGTERM.
// Another signal during the shutdown terminates the process.
func runComponentUntilSignal(ctx context.Context, name string, service component) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)
	return runComponent(ctx, name, service)
}

// runComponent starts the service and shuts it down when ctx is done or the service fails to start.
func runComponent(ctx context.Context, name string, service component) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	startErr := make(chan error, 1)
	go func() {
		err := service.Start(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			startErr <- err
		}
	}()
	slog.InfoContext(ctx, "NeoShowcase service started", "service", name)

	var failed error
	select {
	case <-ctx.Done():
	case failed = <-startErr:
		slog.ErrorContext(ctx, "failed to start service", "error", failed)
	}
	cancel()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	return errors.Join(failed, service.Shutdown(shutdownCtx))
}

func main() {
	// Initialize OpenTelemetry tracer provider for trace ID generation
	if err := observability.InitTracerProvider("neoshowcase"); err != nil {
		slog.Error("failed to initialize tracer provider", "error", err)
		os.Exit(1)
	}

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
		slog.Error("failed to exec", "error", err)
		os.Exit(1)
	}
}
