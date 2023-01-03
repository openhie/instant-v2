package commands

import (
	"cli/cmd/pkg"
	"cli/cmd/project"

	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	// TODO: add commands
	cmd.AddCommand(
		pkg.DeclarePackageCommand(),
		project.ProjectGenerateCommand(),
	)
}
