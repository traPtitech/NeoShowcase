package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// PrintVersion バージョンを出力する
func PrintVersion(cmd *cobra.Command, _ []string) {
	cmd = cmd.Root()
	log.Infof("%s - %s", cmd.Use, cmd.Version)
}
