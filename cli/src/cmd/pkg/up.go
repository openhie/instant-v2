package pkg

import (
	"context"

	"cli/cmd/flags"
	"cli/core/deploy"
	"cli/core/parse"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Stand a package back up after it has been brought down",
		Run: func(cmd *cobra.Command, args []string) {
			packageSpec, config, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(context.Background(), err)
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

	return cmd
}
