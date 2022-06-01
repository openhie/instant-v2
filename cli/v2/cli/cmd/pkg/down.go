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
			name, err := cmd.Flags().GetStringSlice("name")
			util.LogError(err)
			fmt.Printf("Init %s", name)
		},
	}

	flags := cmd.Flags()

	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")

	return cmd
}
