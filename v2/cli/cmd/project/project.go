package project

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func handleProjectAction(init, up, down, remove bool) {
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

func DeclareProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project level commands",
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

				handleProjectAction(init, up, down, remove)
			}
		},
	}
	flags := cmd.Flags()

	flags.BoolP("init", "i", false, "Initialize all packages in the project")
	flags.BoolP("up", "u", false, "Up all packages in the project")
	flags.BoolP("down", "d", false, "Down all packages in the project")
	flags.BoolP("remove", "r", false, "Remove all packages in the project")

	cmd.MarkFlagsMutuallyExclusive("init", "up", "down", "remove")

	cmd.AddCommand(
		ProjectGenerateCommand(),
	)

	return cmd
}
