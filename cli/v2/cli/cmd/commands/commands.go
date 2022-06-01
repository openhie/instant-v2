package commands

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/config"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/pkg"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/project"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/stack"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		config.InitConfigCommand(),
		pkg.InitPackageCommand(),
		project.InitProjectCommand(),
		stack.InitStackCommand(),
	)
}
