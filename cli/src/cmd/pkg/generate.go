package pkg

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"cli/core"
	prompt "cli/prompt/package"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

// TODO(MarkL): Write tests for this once this functionality is introduced
func PackageGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate a new package",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			ex, err := os.Executable()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			pwd := filepath.Dir(ex)

			resp, err := prompt.GeneratePackagePrompt()
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			packagePath := path.Join(pwd, resp.Id)
			err = os.Mkdir(packagePath, os.ModePerm)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}

			generatePackageSpec := core.GeneratePackageSpec{
				Id:             resp.Id,
				Name:           resp.Name,
				Stack:          resp.Stack,
				Description:    resp.Description,
				Type:           resp.Type,
				IncludeDevFile: resp.IncludeDevFile,
			}
			err = core.GeneratePackage(packagePath, generatePackageSpec)
			if err != nil {
				log.Error(ctx, err)
				panic(err)
			}
		},
	}

	return cmd
}
