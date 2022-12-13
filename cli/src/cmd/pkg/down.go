package pkg

import (
	"context"

	"cli/cmd/flags"
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
