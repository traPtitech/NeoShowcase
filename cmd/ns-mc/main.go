package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/traPtitech/neoshowcase/pkg/domain/web"
	"github.com/traPtitech/neoshowcase/pkg/interface/handler"
	"github.com/traPtitech/neoshowcase/pkg/usecase"
	"github.com/traPtitech/neoshowcase/pkg/util/cli"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

var (
	port           int
	pubkeyFilePath string
	cookieName     string
)

var rootCommand = &cobra.Command{
	Use:              "ns-mc",
	Short:            "NeoShowcase MemberCheckerServer",
	Version:          fmt.Sprintf("%s (%s)", version, revision),
	PersistentPreRun: cli.PrintVersion,
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
					log.Fatalf("failed to start server: %+v", err)
				}
			}()

			cli.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	}
	return cmd
}

func main() {
	rootCommand.AddCommand(
		serveCommand(),
	)
	flags := rootCommand.PersistentFlags()
	cli.SetupDebugFlag(flags)
	cli.SetupLogLevelFlag(flags)

	flags.IntVarP(&port, "port", "p", cli.GetIntEnvOrDefault("NS_MC_PORT", 8080), "port num")
	flags.StringVarP(&pubkeyFilePath, "pubkey-file", "k", cli.GetEnvOrDefault("NS_MC_PUBKEY_FILE", ""), "public key PEM file path")
	flags.StringVarP(&cookieName, "cookie-name", "c", cli.GetEnvOrDefault("NS_MC_COOKIE_NAME", "traP_ext_token"), "token cookie name")

	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("failed to exec: %+v", err)
	}
}

func providePubKeyPEM() (usecase.TrapShowcaseJWTPublicKeyPEM, error) {
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

	if len(pubkeyFilePath) > 0 {
		b, err := os.ReadFile(pubkeyFilePath)
		if err != nil {
			return "", err
		}
		return usecase.TrapShowcaseJWTPublicKeyPEM(b), nil
	}
	return defaultPublicKeyPEM, nil
}

func provideServerConfig(h handler.MemberCheckHandler) web.Config {
	return web.Config{
		Port: port,
		SetupRoute: func(e *echo.Echo) {
			e.Any("/*", web.UnwrapHandler(h))
		},
	}
}

func provideTokenCookieName() handler.TokenCookieName {
	return handler.TokenCookieName(cookieName)
}
