package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version  string
	revision string
)

func SetVersion(ver, rev string) {
	version = ver
	revision = rev
}

func GetVersion() (ver, rev string) {
	return version, revision
}

// PrintVersion バージョンを出力する
func PrintVersion(cmd *cobra.Command, _ []string) {
	cmd = cmd.Root()
	log.Infof("%s - %s", cmd.Use, cmd.Version)
}
