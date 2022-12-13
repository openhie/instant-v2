package pkg

import (
	"cli/cmd/flags"
	"cli/core/parse"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Stand a package back up after it has been brought down",
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
