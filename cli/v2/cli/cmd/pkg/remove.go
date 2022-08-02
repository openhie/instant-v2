package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitRemoveCommand() *cobra.Command {
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

			config := core.LoadConfig("config.yml")

			core.LaunchPackage(packageSpec, config)
		},
	}

	flags := cmd.Flags()

	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("only", false, "Only remove the package(s) provided and not their dependency packages")

	return cmd
}
