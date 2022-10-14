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

			config, err := getConfigFromParams(cmd)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			packageSpec, err := getPackageSpecFromParams(cmd, config)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			packageSpec.DeployCommand = "destroy"
			packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			err = validate(cmd, config)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			err = core.LaunchPackage(*packageSpec, *config)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
		},
	}

	setPackageActionFlags(cmd)

	return cmd
}
