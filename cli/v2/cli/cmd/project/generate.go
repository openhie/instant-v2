package project

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func generateProject(name string) {
	fmt.Println("To be implemented")
	fmt.Println(name)
}

func InitGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project",
		Run: func(cmd *cobra.Command, args []string) {
			name := util.GetFlagOrDefaultString(cmd, "name")
			generateProject(name)
		},
	}

	flags := cmd.Flags()

	flags.String("name", "project", "The name of the new project")

	return cmd
}
