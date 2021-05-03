package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/appdefaults"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"github.com/traPtitech/neoshowcase/pkg/domain/builder"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/admindb"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/storage"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
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
	Use:              "ns-builder",
	Short:            "NeoShowcase BuilderService",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cliutil.PrintVersion,
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := New(c)
			if err != nil {
				return err
			}

			go func() {
				err := service.Start(context.Background())
				if err != nil {
					log.Fatal(err)
				}
			}()

			log.Info("NeoShowcase BuilderService started")
			cliutil.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return service.Shutdown(ctx)
		},
	}
	return cmd
}

func main() {
	rand.Seed(time.Now().UnixNano())
	cobra.OnInitialize(cliutil.CobraOnInitializeFunc(&configFilePath, "NS_BUILDER", &c))

	rootCommand.AddCommand(
		runCommand(),
		cliutil.PrintConfCommand(&c),
	)

	flags := rootCommand.PersistentFlags()
	flags.StringVarP(&configFilePath, "config", "c", "", "config file path")
	cliutil.SetupDebugFlag(flags)
	cliutil.SetupLogLevelFlag(flags)

	viper.SetDefault("buildkit.address", appdefaults.Address)
	viper.SetDefault("buildkit.registry", "localhost:5000")
	viper.SetDefault("grpc.port", 10000)
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local.dir", "/neoshowcase")
	viper.SetDefault("storage.s3.bucket", "neoshowcase")
	viper.SetDefault("storage.s3.endpoint", "")
	viper.SetDefault("storage.s3.region", "")
	viper.SetDefault("storage.s3.accessKey", "")
	viper.SetDefault("storage.s3.accessSecret", "")
	viper.SetDefault("storage.swift.username", "")
	viper.SetDefault("storage.swift.apiKey", "")
	viper.SetDefault("storage.swift.tenantName", "")
	viper.SetDefault("storage.swift.tenantId", "")
	viper.SetDefault("storage.swift.container", "neoshowcase")
	viper.SetDefault("storage.swift.authUrl", "")

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

func provideAdminDBConfig(c Config) admindb.Config {
	return c.DB
}

func provideStorageConfig(c Config) storage.Config {
	return c.Storage
}

func provideGRPCPort(c Config) grpc.TCPListenPort {
	return grpc.TCPListenPort(c.GRPC.Port)
}

func provideDockerImageRegistry(c Config) builder.DockerImageRegistryString {
	return builder.DockerImageRegistryString(c.Buildkit.Registry)
}

func initStorage(c storage.Config) (storage.Storage, error) {
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

func initBuildkitClient(c Config) (*buildkit.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildkit.New(ctx, c.Buildkit.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Buildkit Client: %w", err)
	}
	return client, nil
}