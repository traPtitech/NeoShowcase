package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/traPtitech/neoshowcase/pkg/domain"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/builtin"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/staticserver/caddy"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/usecase/healthcheck"
	"github.com/traPtitech/neoshowcase/pkg/usecase/ssgen"
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
	Use:              "ns-ssgen",
	Short:            "NeoShowcase StaticSiteGenerator",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cli.PrintVersion,
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := New(c)
			if err != nil {
				return err
			}

			go func() {
				err := server.Start(context.Background())
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("failed to start server: %+v", err)
				}
			}()

			log.Info("NeoShowcase StaticSiteGenerator started")
			cli.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	}
	return cmd
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

func provideHealthCheckFunc(gen ssgen.GeneratorService) healthcheck.Func {
	return gen.Healthy
}

func provideStaticServer(c Config) (domain.StaticServer, error) {
	switch c.Server.Type {
	case "builtIn":
		return builtin.NewServer(c.Server.BuiltIn, c.ArtifactsRoot), nil
	case "caddy":
		return caddy.NewServer(c.Server.Caddy), nil
	default:
		return nil, errors.Errorf("invalid static server type: %v", c.Server.Type)
	}
}

func provideStaticServerDocumentRootPath(c Config) domain.StaticServerDocumentRootPath {
	return domain.StaticServerDocumentRootPath(c.ArtifactsRoot)
}

func provideStorage(c domain.StorageConfig) (domain.Storage, error) {
	switch strings.ToLower(c.Type) {
	case "local":
		return storage.NewLocalStorage(c.Local.Dir)
	case "s3":
		return storage.NewS3Storage(c.S3.Bucket, c.S3.AccessKey, c.S3.AccessSecret, c.S3.Region, c.S3.Endpoint)
	case "swift":
		return storage.NewSwiftStorage(c.Swift.Container, c.Swift.UserName, c.Swift.APIKey, c.Swift.TenantName, c.Swift.TenantID, c.Swift.AuthURL)
	default:
		return nil, fmt.Errorf("unknown storage: %s", c.Type)
	}
}
