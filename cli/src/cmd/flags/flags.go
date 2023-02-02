package flags

import (
	"cli/core/state"

	"github.com/spf13/cobra"
)

func SetPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	setCommonActionFlags(cmd)

	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.StringP("profile", "p", "", "The profile name to load parameters from (defined in config.yml)")
}

// disables certain flags for project level commands, but allows for usage
// of parseAndPrepareLaunch function without receiving errors
func SetProjectActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	// allowed flags
	setCommonActionFlags(cmd)

	// disabled flags
	flags.StringSlice("name", nil, "")
	flags.String("profile", "", "")

	cmd.Flags().MarkHidden("name")
	cmd.Flags().MarkHidden("profile")
}

func setCommonActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")
	flags.BoolP("dev", "d", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.BoolP("only", "o", false, "Ignore package dependencies")
	flags.StringSliceVar(&state.EnvFiles, "env-file", nil, "env file")
	flags.StringVar(&state.ConfigFile, "config", "", "config file (default is $WORKING_DIR/config.yaml)")
	flags.StringSliceP("env-var", "e", nil, "Env var(s) to set or overwrite")
}
