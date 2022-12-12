package pkg

import (
	"github.com/spf13/cobra"
)

func packageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	setPackageActionFlags(cmd)

	return cmd
}
