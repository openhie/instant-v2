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
			config, err := getConfigFromParams(cmd)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			packageSpec, err := getPackageSpecFromParams(cmd, config)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			err = validate(cmd, config)
			if err != nil {
				log.Error(context.Background(), err)
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
				log.Error(context.Background(), err)
				panic(err)
			}
		},
	}

	setPackageActionFlags(cmd)
	return cmd
}
