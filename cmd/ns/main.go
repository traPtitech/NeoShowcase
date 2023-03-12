package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/grpc"
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

			// TODO: context
			go func() {
				err := service.Start(context.Background())
				if err != nil {
					log.Fatal(err)
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
	viper.SetDefault("image.registry", "localhost")
	viper.SetDefault("image.namePrefix", "ns-apps/")
	viper.SetDefault("builder.addr", "")
	viper.SetDefault("builder.insecure", false)
	viper.SetDefault("ssgen.addr", "")
	viper.SetDefault("ssgen.insecure", false)
	viper.SetDefault("grpc.port", 5000)
	viper.SetDefault("http.port", 10000)
	viper.SetDefault("http.debug", false)
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

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

var handlerSet = wire.NewSet()

type Router struct {
}

func (r *Router) SetupRoute(e *echo.Echo) {
	_ = e.Group("")
}

func provideGRPCPort(c Config) grpc.TCPListenPort {
	return grpc.TCPListenPort(c.GRPC.Port)
}

func provideWebServerConfig(router web.Router) web.Config {
	return web.Config{
		Port:   c.HTTP.Port,
		Router: router,
	}
}
