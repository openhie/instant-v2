package pkg

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			name := util.GetFlagOrDefaultString(cmd, "name")
			fmt.Printf("Init %s", name)
		},
	}

	flags := cmd.Flags()

	flags.StringP("name", "n", "package", "The name(s) of the package(s)")

	return cmd
}
