package commands

import (
	"cli/cmd/pkg"

	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	// TODO: add commands
	cmd.AddCommand(pkg.DeclarePackageCommand())
}
