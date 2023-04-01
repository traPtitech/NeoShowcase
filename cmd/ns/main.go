package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Use:              "ns",
	Short:            "NeoShowcase Core API Server",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cli.PrintVersion,
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := New(c)
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

			log.Info("NeoShowcase ApiServer started")
			cli.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return service.Shutdown(ctx)
		},
	}
	return cmd
}

func main() {
	rand.Seed(time.Now().UnixNano())
	cobra.OnInitialize(cli.CobraOnInitializeFunc(&configFilePath, "NS_APISERVER", &c))

	rootCommand.AddCommand(
		runCommand(),
		cli.PrintConfCommand(&c),
	)

	flags := rootCommand.PersistentFlags()
	flags.StringVarP(&configFilePath, "config", "c", "", "config file path")
	cli.SetupDebugFlag(flags)
	cli.SetupLogLevelFlag(flags)

	viper.SetDefault("debug", false)
	viper.SetDefault("mode", "docker")
	viper.SetDefault("repository.cacheDir", "")
	viper.SetDefault("repository.privateKeyFile", "")
	viper.SetDefault("image.registry.scheme", "https")
	viper.SetDefault("image.registry.addr", "localhost")
	viper.SetDefault("image.registry.username", "")
	viper.SetDefault("image.registry.password", "")
	viper.SetDefault("image.namePrefix", "ns-apps/")
	viper.SetDefault("ss.service.namespace", "default")
	viper.SetDefault("ss.service.kind", "Service")
	viper.SetDefault("ss.service.name", "")
	viper.SetDefault("ss.service.port", 80)
	viper.SetDefault("ss.url", "")
	viper.SetDefault("docker.confdir", "/opt/traefik/conf")
	viper.SetDefault("grpc.app.port", 5000)
	viper.SetDefault("grpc.component.port", 10000)
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.database", "neoshowcase")
	viper.SetDefault("db.connection.maxOpen", 0)
	viper.SetDefault("db.connection.maxIdle", 2)
	viper.SetDefault("db.connection.lifetime", 0)
	viper.SetDefault("mariadb.host", "127.0.0.1")
	viper.SetDefault("mariadb.port", 3306)
	viper.SetDefault("mariadb.adminUser", "root")
	viper.SetDefault("mariadb.adminPassword", "password")
	viper.SetDefault("mongodb.host", "127.0.0.1")
	viper.SetDefault("mongodb.port", 27017)
	viper.SetDefault("mongodb.adminUser", "root")
	viper.SetDefault("mongodb.adminPassword", "password")
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
		log.Fatalf("failed to exec: %+v", err)
	}
}
