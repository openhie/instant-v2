package pkg

import (
	"cli/cmd/flags"

	"github.com/spf13/cobra"
)

func packageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	flags.SetPackageActionFlags(cmd)

	return cmd
}
