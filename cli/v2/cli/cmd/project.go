package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func handleGenerate() {
	fmt.Println("GENERATE")
}

func handleAction(init, up, down, remove bool) {
	switch true {
	case init:
		fmt.Println("INIT")
	case up:
		fmt.Println("UP")
	case down:
		fmt.Println("DOWN")
	case remove:
		fmt.Println("REMOVE")
	default:
		panic("Invalid action entered")
	}
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project level commands",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("generate").Changed {
			handleGenerate()
		}

		if cmd.Flag("init").Changed || cmd.Flag("up").Changed || cmd.Flag("down").Changed || cmd.Flag("remove").Changed {
			init, err := cmd.Flags().GetBool("init")
			logError(err)
			up, err := cmd.Flags().GetBool("up")
			logError(err)
			down, err := cmd.Flags().GetBool("down")
			logError(err)
			remove, err := cmd.Flags().GetBool("remove")
			logError(err)

			handleAction(init, up, down, remove)
		}
		fmt.Println("project called")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.Flags().StringP("generate", "g", "project", "Generate a new project")
	projectCmd.Flags().BoolP("init", "i", false, "Initialize all packages in the project")
	projectCmd.Flags().BoolP("up", "u", false, "Up all packages in the project")
	projectCmd.Flags().BoolP("down", "d", false, "Down all packages in the project")
	projectCmd.Flags().BoolP("remove", "r", false, "Remove all packages in the project")
	projectCmd.MarkFlagsMutuallyExclusive("init", "up", "down", "remove", "generate")

}
