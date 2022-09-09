package pkg

import (
	"os"
	"path"
	"path/filepath"

	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	prompt "github.com/openhie/package-starter-kit/cli/v2/cli/prompt/package"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
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
