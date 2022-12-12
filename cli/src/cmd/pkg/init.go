package pkg

import (
	"github.com/spf13/cobra"
)

func packageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	setPackageActionFlags(cmd)

	return cmd
}
