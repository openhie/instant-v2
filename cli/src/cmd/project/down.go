package project

import (
	"context"

	pFlags "cli/cmd/flags"
	"cli/core/deploy"
	"cli/core/parse"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func projectDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Down all packages in the project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			err := checkInvalidFlags(cmd)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			packageSpec, config, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
			packageSpec.Packages = config.Packages
			packageSpec.CustomPackages = config.CustomPackages

			err = deploy.LaunchDeploymentContainer(packageSpec, config)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
		},
	}

	pFlags.SetProjectActionFlags(cmd)

	return cmd
}
