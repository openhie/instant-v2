package pkg

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "up",
		Aliases: []string{"i"},
		Short:   "Stand a package back up after it has been brought down",
		Run: func(cmd *cobra.Command, args []string) {
			name := util.GetFlagOrDefaultString(cmd, "name")
			fmt.Printf("Up %s", name)
		},
	}

	flags := cmd.Flags()

	flags.StringP("name", "n", "package", "The name(s) of the package(s)")

	return cmd
}
