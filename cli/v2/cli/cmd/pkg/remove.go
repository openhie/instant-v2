package pkg

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"i"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			name := util.GetFlagOrDefaultString(cmd, "name")
			fmt.Printf("Remove %s", name)
		},
	}

	flags := cmd.Flags()

	flags.StringP("name", "n", "package", "The name(s) of the package(s)")

	return cmd
}
