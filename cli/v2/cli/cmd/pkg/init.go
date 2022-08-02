package pkg

import (
	"fmt"
	"strings"

	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageInitCommand(global *types.Global) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			packageNames, err := cmd.Flags().GetStringSlice("name")
			util.LogError(err)
			isDev, err := cmd.Flags().GetBool("dev")
			util.LogError(err)
			isOnly, err := cmd.Flags().GetBool("only")
			util.LogError(err)

			var envVariables []string
			allEnvVars := global.EnvVarViper.AllSettings()
			for key, element := range allEnvVars {
				envVariables = append(envVariables, fmt.Sprintf("%v=%v", strings.ToUpper(key), element))
			}

			packageSpec := core.PackageSpec{
				Packages:             packageNames,
				EnvironmentVariables: envVariables,
				IsDev:                isDev,
				IsOnly:               isOnly,
				DeployCommand:        cmd.Use,
			}

			var config core.Config
			err = global.ConfigViper.Unmarshal(&config)
			util.PanicError(err)

			core.LaunchPackage(packageSpec, config)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("dev", false, "Should the package launch in dev mode")
	flags.Bool("only", false, "Should the package launch without launching it's dependency packages")

	return cmd
}
