package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/spf13/cobra"
)

func DeclarePackageCommand(global *types.Global) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	cmd.AddCommand(
		PackageInitCommand(global),
		PackageUpCommand(),
		PackageDownCommand(),
		PackageRemoveCommand(global),
	)

	return cmd
}
