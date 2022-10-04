package commands

import (
	"cli/cmd/config"
	"cli/cmd/pkg"
	"cli/cmd/project"
	"cli/cmd/stack"

	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		config.DeclareConfigCommand(),
		pkg.DeclarePackageCommand(),
		project.DeclareProjectCommand(),
		stack.DeclareStackCommand(),
	)
}
