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

func packageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "destroy"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			packageSpec, config, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}
			packageSpec.DeployCommand = "destroy"

			if len(packageSpec.Packages) < 1 {
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
