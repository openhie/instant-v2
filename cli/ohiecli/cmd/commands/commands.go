package commands

import (
	"ohiecli/cmd/config"
	"ohiecli/cmd/pkg"
	"ohiecli/cmd/project"
	"ohiecli/cmd/stack"

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
