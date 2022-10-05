package project

import (
	"cli/core"
	prompt "cli/prompt/project"
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func ProjectGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			resp, err := prompt.GenerateProjectPrompt()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			config := core.Config{
				Image:         resp.ProjectImage,
				ProjectName:   resp.ProjectName,
				PlatformImage: resp.PlatformImage,
			}
			err = core.GenerateConfigFile(&config)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}
		},
	}

	return cmd
}
