package pkg

import (
	"cli/cmd/flags"

	"github.com/spf13/cobra"
)

func packageUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"u"},
		Short:   "Stand a package back up after it has been brought down",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	flags.SetPackageActionFlags(cmd)

	return cmd
}
