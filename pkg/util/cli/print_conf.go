package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
)

// PrintConfCommand 設定確認・ベース設定プリントコマンド
func PrintConfCommand(config interface{}) *cobra.Command {
	return &cobra.Command{
		Use:   "print-conf",
		Short: "Print loaded config variables",
		Run: func(cmd *cobra.Command, args []string) {
			if err := yaml.NewEncoder(os.Stdout).Encode(config); err != nil {
				log.Fatalf("unable to marshal config to YAML: %v", err)
			}
		},
	}
}
