package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := getConfigFromParams(cmd)
			util.PanicError(err)
			packageSpec, err := getPackageSpecFromParams(cmd)
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
