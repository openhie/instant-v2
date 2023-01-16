package flags

import "github.com/spf13/cobra"

func SetPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.BoolP("dev", "d", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.BoolP("only", "o", false, "Ignore package dependencies")
	flags.StringP("profile", "p", "", "The profile name to load parameters from (defined in config.yml)")
	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")
}

// disables certain flags for project level commands, but allows for usage
// of parseAndPrepareLaunch function without receiving errors
func SetProjectActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	// allowed flags
	flags.BoolP("dev", "d", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.BoolP("only", "o", false, "Ignore package dependencies")
	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")

	// disabled flags
	flags.StringSlice("name", nil, "")
	flags.String("profile", "", "")

	cmd.Flags().MarkHidden("name")
	cmd.Flags().MarkHidden("profile")
}
