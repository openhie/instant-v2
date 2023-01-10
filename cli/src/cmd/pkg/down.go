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

func packageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run: func(cmd *cobra.Command, args []string) {
			packageSpec, config, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

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
