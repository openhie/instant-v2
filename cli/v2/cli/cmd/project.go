package cmd

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/pkg"
	"github.com/spf13/cobra"
)

func handleGenerate() {
	fmt.Println("GENERATE")
	fmt.Println("To be implemented")
}

func handleAction(init, up, down, remove bool) {
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
			pkg.LogError(err)
			up, err := cmd.Flags().GetBool("up")
			pkg.LogError(err)
			down, err := cmd.Flags().GetBool("down")
			pkg.LogError(err)
			remove, err := cmd.Flags().GetBool("remove")
			pkg.LogError(err)

			handleAction(init, up, down, remove)
		}
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.Flags().BoolP("init", "i", false, "Initialize all packages in the project")
	projectCmd.Flags().BoolP("up", "u", false, "Up all packages in the project")
	projectCmd.Flags().BoolP("down", "d", false, "Down all packages in the project")
	projectCmd.Flags().BoolP("remove", "r", false, "Remove all packages in the project")
	projectCmd.MarkFlagsMutuallyExclusive("init", "up", "down", "remove")

}
