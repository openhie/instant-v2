package pkg

import (
	"fmt"
	"strings"

	viperUtil "github.com/openhie/package-starter-kit/cli/v2/cli/cmd/util"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "down",
		Aliases: []string{"d"},
		Short:   "Bring a package down without removing volumes or configs",
		Run: func(cmd *cobra.Command, args []string) {
			packageNames, err := cmd.Flags().GetStringSlice("name")
			util.LogError(err)

			envFiles, err := cmd.Flags().GetStringSlice("env-file")
			util.LogError(err)
			envViper := viperUtil.GetEnvironmentVariableViper(envFiles)
			var envVariables []string
			allEnvVars := envViper.AllSettings()
			for key, element := range allEnvVars {
				envVariables = append(envVariables, fmt.Sprintf("%v=%v", strings.ToUpper(key), element))
			}

			configFile, err := cmd.Flags().GetString("config")
			util.LogError(err)
			configViper := viperUtil.GetConfigViper(configFile)
			var config core.Config
			err = configViper.Unmarshal(&config)
			util.PanicError(err)

			packageSpec := core.PackageSpec{
				Packages:             packageNames,
				DeployCommand:        cmd.Use,
				EnvironmentVariables: envVariables,
			}

			core.LaunchPackage(packageSpec, config)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")

	return cmd
}
