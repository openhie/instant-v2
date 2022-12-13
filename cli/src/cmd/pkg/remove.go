package pkg

import (
	"cli/cmd/flags"
	"cli/core/parse"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "destroy"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
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
