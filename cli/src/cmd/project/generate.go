package project

import (
	"context"

	"cli/core/generate"
	"cli/core/prompt"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func projectGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			config, err := prompt.GenerateProjectPrompt()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			err = generate.GenerateConfigFile(&config)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}
		},
	}

	return cmd
}
