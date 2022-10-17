package pkg

import (
	"cli/core"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func PackageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
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

			packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			for _, pack := range packageSpec.Packages {
				for _, customPack := range config.CustomPackages {
					if pack == customPack.Id {
						packageSpec.CustomPackages = append(packageSpec.CustomPackages, customPack)
					}
				}
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
