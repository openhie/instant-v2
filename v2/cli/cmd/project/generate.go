package project

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	prompt "github.com/openhie/package-starter-kit/cli/v2/cli/prompt/project"
	"github.com/spf13/cobra"
)

func ProjectGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project",
		Run: func(cmd *cobra.Command, args []string) {
			resp := prompt.GenerateProjectPrompt()
			config := core.Config{
				Image:         resp.ProjectImage,
				ProjectName:   resp.ProjectName,
				PlatformImage: resp.PlatformImage,
			}
			core.GenerateConfigFile(&config)
		},
	}

	return cmd
}
