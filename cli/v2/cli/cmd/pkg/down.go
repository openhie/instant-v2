package pkg

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run: func(cmd *cobra.Command, args []string) {
			name := util.GetFlagOrDefaultString(cmd, "name")
			fmt.Printf("Down %s", name)
		},
	}

	flags := cmd.Flags()

	flags.StringP("name", "n", "package", "The name(s) of the package(s)")

	return cmd
}
