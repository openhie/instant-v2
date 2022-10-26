package pkg

import (
	"cli/core"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func PackageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			packageSpec, config, err := packageActionHook(cmd, []string{})
			if err != nil {
				log.Error(ctx, err)
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
