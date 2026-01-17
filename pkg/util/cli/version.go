package cli

import (
	"log/slog"

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
	slog.Info("neoshowcase", "name", cmd.Use, "version", cmd.Version)
}
