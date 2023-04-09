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
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	cobra.OnInitialize(cli.CobraOnInitializeFunc(&configFilePath, "NS", &c))

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

	viper.SetDefault("docker.confDir", "/opt/traefik/conf")
	viper.SetDefault("docker.ss.url", "")
	viper.SetDefault("docker.network", "neoshowcase_apps")
	viper.SetDefault("docker.certResolver", "nsresolver")

	viper.SetDefault("k8s.ss.namespace", "default")
	viper.SetDefault("k8s.ss.kind", "Service")
	viper.SetDefault("k8s.ss.name", "")
	viper.SetDefault("k8s.ss.port", 80)
	viper.SetDefault("k8s.namespace", "neoshowcase-apps")
	viper.SetDefault("k8s.tls.type", "traefik")
	viper.SetDefault("k8s.tls.traefik.certResolver", "nsresolver")
	viper.SetDefault("k8s.tls.certManager.issuer.name", "cert-issuer")
	viper.SetDefault("k8s.tls.certManager.issuer.kind", "ClusterIssuer")
	viper.SetDefault("k8s.imagePullSecret", "")

	viper.SetDefault("log.type", "loki")
	viper.SetDefault("log.loki.endpoint", "http://loki:3100")
	viper.SetDefault("log.loki.appIDLabel", "neoshowcase_trap_jp_appId")

	viper.SetDefault("web.app.port", 5000)
	viper.SetDefault("web.component.port", 10000)

	viper.SetDefault("repository.privateKeyFile", "")

	viper.SetDefault("image.registry.scheme", "https")
	viper.SetDefault("image.registry.addr", "localhost")
	viper.SetDefault("image.registry.username", "")
	viper.SetDefault("image.registry.password", "")
	viper.SetDefault("image.namePrefix", "ns-apps/")

	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to exec: %+v", err)
	}
}
