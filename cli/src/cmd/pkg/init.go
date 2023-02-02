package pkg

import (
	"context"

	"cli/cmd/completion"
	"cli/cmd/flags"
	"cli/core/deploy"
	"cli/core/parse"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			packageSpec, config, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			if len(packageSpec.Packages) < 1 && len(packageSpec.CustomPackages) < 1 {
				log.Error(context.Background(), ErrNoPackages)
				panic(err)
			}

			err = deploy.LaunchDeploymentContainer(packageSpec, config)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}
		},
	}

	flags.SetPackageActionFlags(cmd)
	completion.FlagCompletion(cmd)

	return cmd
}
