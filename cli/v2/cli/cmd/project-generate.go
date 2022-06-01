package cmd

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/pkg"
	"github.com/spf13/cobra"
)

func generateProject(name string) {
	fmt.Println("generate called")
	fmt.Println(name)
}

var projectGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new project",
	Run: func(cmd *cobra.Command, args []string) {
		name := pkg.GetFlagOrDefaultString(cmd, "name")
		generateProject(name)
	},
}

func init() {
	projectCmd.AddCommand(projectGenerateCmd)

	projectGenerateCmd.Flags().String("name", "project", "The name of the new project")
}
