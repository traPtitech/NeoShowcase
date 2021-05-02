package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"github.com/traPtitech/neoshowcase/pkg/infrastructure/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
)

const defaultPublicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAraewUw7V1hiuSgUvkly9
X+tcIh0e/KKqeFnAo8WR3ez2tA0fGwM+P8sYKHIDQFX7ER0c+ecTiKpo/Zt/a6AO
gB/zHb8L4TWMr2G4q79S1gNw465/SEaGKR8hRkdnxJ6LXdDEhgrH2ZwIPzE0EVO1
eFrDms1jS3/QEyZCJ72oYbAErI85qJDF/y/iRgl04XBK6GLIW11gpf8KRRAh4vuh
g5/YhsWUdcX+uDVthEEEGOikSacKZMFGZNi8X8YVnRyWLf24QTJnTHEv+0EStNrH
HnxCPX0m79p7tBfFC2ha2OYfOtA+94ZfpZXUi2r6gJZ+dq9FWYyA0DkiYPUq9QMb
OQIDAQAB
-----END PUBLIC KEY-----
`

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	port           int
	pubkeyFilePath string
)

var rootCommand = &cobra.Command{
	Use:              "ns-mc",
	Short:            "NeoShowcase MemberCheckerServer",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cliutil.PrintVersion,
}

func serveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := NewServer()
			if err != nil {
				return err
			}

			go func() {
				err := server.Start(context.Background())
				if err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()

			cliutil.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	}
	return cmd
}

func init() {
	rand.Seed(time.Now().UnixNano())
	rootCommand.AddCommand(
		serveCommand(),
	)
	flags := rootCommand.PersistentFlags()
	cliutil.SetupDebugFlag(flags)
	cliutil.SetupLogLevelFlag(flags)

	flags.IntVarP(&port, "port", "p", cliutil.GetIntEnvOrDefault("NS_MC_PORT", 8081), "port num")
	flags.StringVarP(&pubkeyFilePath, "pubkey-file", "k", cliutil.GetEnvOrDefault("NS_MC_PUBKEY_FILE", ""), "public key PEM file path")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}

var handlerSet = wire.NewSet(
	handler.NewMemberCheckHandler,
)

type Router struct {
	h handler.MemberCheckHandler
}

func (r *Router) SetupRoute(e *echo.Echo) {
	e.GET("/", web.UnwrapHandler(r.h))
}

func providePubKeyPEM() (usecase.TrapShowcaseJWTPublicKeyPEM, error) {
	if len(pubkeyFilePath) > 0 {
		b, err := os.ReadFile(pubkeyFilePath)
		if err != nil {
			return "", err
		}
		return usecase.TrapShowcaseJWTPublicKeyPEM(b), nil
	}
	return defaultPublicKeyPEM, nil
}

func provideServerConfig(router web.Router) web.Config {
	return web.Config{
		Port:   port,
		Router: router,
	}
}
