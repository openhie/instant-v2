package pkg

import (
	"cli/cmd/flags"

	"github.com/spf13/cobra"
)

func packageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	flags.SetPackageActionFlags(cmd)

	return cmd
}
