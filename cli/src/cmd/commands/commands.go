package commands

import (
	"cli/cmd/completion"
	"cli/cmd/pkg"
	"cli/cmd/project"
	"cli/cmd/version"

	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		pkg.DeclarePackageCommand(),
		project.DeclareProjectCommand(),
		completion.GenCompletionCommand(),
		version.VersionCommand(),
	)
}
