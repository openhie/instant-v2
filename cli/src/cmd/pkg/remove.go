package pkg

import (
	"cli/core"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func PackageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "destroy"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			packageSpec, config, err := parseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
			packageSpec.DeployCommand = "destroy"

			err = core.LaunchDeploymentContainer(*packageSpec, *config)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
		},
	}

	setPackageActionFlags(cmd)

	return cmd
}
