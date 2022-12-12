package pkg

import (
	"github.com/spf13/cobra"
)

func setPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.BoolP("dev", "d", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.BoolP("only", "o", false, "Ignore package dependencies")
	flags.StringP("profile", "p", "", "The profile name to load parameters from (defined in config.yml)")
	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")
}

func DeclarePackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	cmd.AddCommand(
		packageInitCommand(),
		packageUpCommand(),
		packageDownCommand(),
		packageRemoveCommand(),
		packageGenerateCommand(),
	)

	return cmd
}
