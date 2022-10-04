package pkg

import (
	"cli/core"
	"cli/util"

	"github.com/spf13/cobra"
)

func PackageUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Stand a package back up after it has been brought down",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := getConfigFromParams(cmd)
			util.PanicError(err)
			packageSpec, err := getPackageSpecFromParams(cmd, config)
			util.PanicError(err)
			packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
			util.PanicError(err)

			err = core.LaunchPackage(*packageSpec, *config)
			util.PanicError(err)
		},
	}

	setPackageActionFlags(cmd)

	return cmd
}
