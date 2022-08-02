package commands

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/config"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/pkg"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/project"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/stack"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command, global *types.Global) {
	cmd.AddCommand(
		config.DeclareConfigCommand(global),
		pkg.DeclarePackageCommand(global),
		project.DeclareProjectCommand(global),
		stack.DeclareStackCommand(global),
	)
}
