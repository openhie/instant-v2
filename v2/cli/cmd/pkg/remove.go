package pkg

import (
	viperUtil "github.com/openhie/package-starter-kit/cli/v2/cli/cmd/util"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r"},
		Short:   "Remove everything related to a package (volumes, configs, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			packageNames, err := cmd.Flags().GetStringSlice("name")
			util.LogError(err)
			isOnly, err := cmd.Flags().GetBool("only")
			util.LogError(err)

			envFiles, err := cmd.Flags().GetStringSlice("env-file")
			util.LogError(err)
			envViper := viperUtil.GetEnvironmentVariableViper(envFiles)
			envVariables := viperUtil.GetEnvVariableString(envViper)

			configFile, err := cmd.Flags().GetString("config")
			util.LogError(err)
			configViper := viperUtil.GetConfigViper(configFile)
			var config core.Config
			err = configViper.Unmarshal(&config)
			util.PanicError(err)

			packageSpec := core.PackageSpec{
				Packages:             packageNames,
				DeployCommand:        "destroy",
				IsOnly:               isOnly,
				EnvironmentVariables: envVariables,
			}

			core.LaunchPackage(packageSpec, config)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")

	return cmd
}
