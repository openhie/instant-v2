package pkg

import (
	"github.com/spf13/cobra"
)

func packageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "destroy"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	SetPackageActionFlags(cmd)

	return cmd
}
