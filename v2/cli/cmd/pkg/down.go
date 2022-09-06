package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/spf13/cobra"
)

func PackageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run: func(cmd *cobra.Command, args []string) {
			config := getConfigFromParams(cmd)
			packageSpec := getPackageSpecFromParams(cmd)
			packageSpec = loadInProfileParams(cmd, *config, *packageSpec)

			core.LaunchPackage(*packageSpec, *config)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("dev", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.Bool("only", false, "Ignore package dependencies")
	flags.String("profile", "", "The profile name to load parameters from (defined in config.yml)")

	return cmd
}
