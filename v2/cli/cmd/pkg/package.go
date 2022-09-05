package pkg

import (
	"github.com/spf13/cobra"
)

func DeclarePackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	flags := cmd.PersistentFlags()
	flags.Bool("dev", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.Bool("only", false, "Ignore package dependencies")

	cmd.AddCommand(
		PackageInitCommand(),
		PackageUpCommand(),
		PackageDownCommand(),
		PackageRemoveCommand(),
	)

	return cmd
}
