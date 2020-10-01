package cli

import "github.com/spf13/cobra"

var (
	Version  string
	Revision string
)

var rootCommand = &cobra.Command{
	Use: "neoshowcase",
}

func Execute() error {
	return rootCommand.Execute()
}
