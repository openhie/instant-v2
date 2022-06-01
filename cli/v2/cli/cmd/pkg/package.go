package pkg

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func handlePackageAction(init, up, down, remove bool, name []string) {
	for _, n := range name {
		fmt.Println(n)
	}

	switch true {
	case init:
		fmt.Println("INIT")
		fmt.Println("To be implemented")
	case up:
		fmt.Println("UP")
		fmt.Println("To be implemented")
	case down:
		fmt.Println("DOWN")
		fmt.Println("To be implemented")
	case remove:
		fmt.Println("REMOVE")
		fmt.Println("To be implemented")
	default:
		panic("Invalid action entered")
	}
}

func InitPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("init").Changed || cmd.Flag("up").Changed || cmd.Flag("down").Changed || cmd.Flag("remove").Changed {
				init, err := cmd.Flags().GetBool("init")
				util.LogError(err)
				up, err := cmd.Flags().GetBool("up")
				util.LogError(err)
				down, err := cmd.Flags().GetBool("down")
				util.LogError(err)
				remove, err := cmd.Flags().GetBool("remove")
				util.LogError(err)

				var name []string
				if cmd.Flag("name").Changed {
					var err error
					name, err = cmd.Flags().GetStringSlice("name")
					util.LogError(err)

					handlePackageAction(init, up, down, remove, name)
				} else {
					fmt.Println("name must also be provided")
				}

			}
		},
	}

	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s) you wish to act on")

	flags.BoolP("init", "i", false, "Initialize a package with relevant configs, volumes and setup")
	flags.BoolP("up", "u", false, "Stand a package back up after it has been brought down")
	flags.BoolP("down", "d", false, "Bring a package down without removing volumes or configs")
	flags.BoolP("remove", "r", false, "Remove everything related to a package (volumes, configs, etc)")

	cmd.MarkFlagsMutuallyExclusive("init", "up", "down", "remove")

	return cmd
}
