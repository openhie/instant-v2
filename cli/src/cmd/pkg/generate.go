package pkg

import (
	"os"
	"path"
	"path/filepath"

	"cli/core"
	prompt "cli/prompt/package"
	"cli/util"

	"github.com/spf13/cobra"
)

func PackageGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate a new package",
		Run: func(cmd *cobra.Command, args []string) {

			ex, err := os.Executable()
			util.PanicError(err)
			pwd := filepath.Dir(ex)

			resp, err := prompt.GeneratePackagePrompt()
			util.PanicError(err)

			packagePath := path.Join(pwd, resp.Id)
			err = os.Mkdir(packagePath, os.ModePerm)
			util.PanicError(err)

			generatePackageSpec := core.GeneratePackageSpec{
				Id:             resp.Id,
				Name:           resp.Name,
				Stack:          resp.Stack,
				Description:    resp.Description,
				Type:           resp.Type,
				IncludeDevFile: resp.IncludeDevFile,
			}
			err = core.GeneratePackage(packagePath, generatePackageSpec)
			util.PanicError(err)
		},
	}

	return cmd
}
