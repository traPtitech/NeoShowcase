package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/traPtitech/neoshowcase/pkg/cliutil"
	"github.com/traPtitech/neoshowcase/pkg/memberchecker"
)

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
			cfg := memberchecker.Config{
				HTTPPort:        port,
				JWTPublicKeyPEM: "",
			}

			if len(pubkeyFilePath) > 0 {
				b, err := ioutil.ReadFile(pubkeyFilePath)
				if err != nil {
					return err
				}
				cfg.JWTPublicKeyPEM = string(b)
			}

			s, err := memberchecker.New(cfg)
			if err != nil {
				return err
			}

			go func() {
				err := s.Start(context.Background())
				if err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()

			cliutil.WaitSIGINT()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return s.Shutdown(ctx)
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
