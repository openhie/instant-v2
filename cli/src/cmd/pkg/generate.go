package pkg

import (
	"context"
	"os"
	"path"

	"cli/core/generate"
	"cli/core/prompt"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func packageGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate a new package",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			resp, err := prompt.GeneratePackagePrompt()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			// Get the current working directory
			cwd, err := os.Getwd()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			packagePath := path.Join(cwd, resp.Id)
			err = os.Mkdir(packagePath, os.ModePerm)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			err = generate.GeneratePackage(packagePath, resp)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
		},
	}

	return cmd
}
