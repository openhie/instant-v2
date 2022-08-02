package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageRemoveCommand(global *types.Global) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			packageNames, err := cmd.Flags().GetStringSlice("name")
			util.LogError(err)
			isOnly, err := cmd.Flags().GetBool("only")
			util.LogError(err)

			packageSpec := core.PackageSpec{
				Packages:      packageNames,
				DeployCommand: "destroy",
				IsOnly:        isOnly,
			}

			var config core.Config
			err = global.ConfigViper.Unmarshal(&config)
			util.LogError(err)

			core.LaunchPackage(packageSpec, config)
		},
	}

	flags := cmd.Flags()

	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("only", false, "Only remove the package(s) provided and not their dependency packages")
	flags.StringSliceP("env-file", "e", nil, "The path to the env file(s)")

	return cmd
}
