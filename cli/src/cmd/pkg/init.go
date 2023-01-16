package pkg

import (
	"context"

	"cli/core/parse"
	"cli/cmd/flags"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			_, _, err := parse.ParseAndPrepareLaunch(cmd)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

		},
	}

	flags.SetPackageActionFlags(cmd)

	return cmd
}
