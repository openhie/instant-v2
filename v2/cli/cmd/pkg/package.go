package pkg

import (
	viperUtil "github.com/openhie/package-starter-kit/cli/v2/cli/cmd/util"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func getConfigFromParams(cmd *cobra.Command) *core.Config {
	var config core.Config
	configFile, err := cmd.Flags().GetString("config")
	util.LogError(err)
	configViper := viperUtil.GetConfigViper(configFile)
	err = configViper.Unmarshal(&config)
	util.PanicError(err)
	return &config
}

func getPackageSpecFromParams(cmd *cobra.Command) *core.PackageSpec {
	packageNames, err := cmd.Flags().GetStringSlice("name")
	util.LogError(err)
	isDev, err := cmd.Flags().GetBool("dev")
	util.LogError(err)
	isOnly, err := cmd.Flags().GetBool("only")
	util.LogError(err)

	envFiles, err := cmd.Flags().GetStringSlice("env-file")
	util.LogError(err)

	envViper := viperUtil.GetEnvironmentVariableViper(envFiles)
	envVariables := viperUtil.GetEnvVariableString(envViper)

	packageSpec := core.PackageSpec{
		Packages:             packageNames,
		EnvironmentVariables: envVariables,
		IsDev:                isDev,
		IsOnly:               isOnly,
		DeployCommand:        cmd.Use,
	}

	return &packageSpec
}

func loadInProfileParams(cmd *cobra.Command, config core.Config, packageSpec core.PackageSpec) *core.PackageSpec {
	profile := core.Profile{}
	profileName, err := cmd.Flags().GetString("profile")
	util.LogError(err)
	for _, p := range config.Profiles {
		if p.Name == profileName {
			profile = p
			break
		}
	}

	if !cmd.Flags().Changed("dev") && profile.Dev {
		packageSpec.IsDev = profile.Dev
	}
	if !cmd.Flags().Changed("only") && profile.Only {
		packageSpec.IsOnly = profile.Only
	}
	if len(profile.Packages) > 0 {
		packageSpec.Packages = append(profile.Packages, packageSpec.Packages...)
	}
	if len(profile.EnvFiles) > 0 {
		envViper := viperUtil.GetEnvironmentVariableViper(profile.EnvFiles)
		envVariables := viperUtil.GetEnvVariableString(envViper)
		packageSpec.EnvironmentVariables = append(envVariables, packageSpec.EnvironmentVariables...)
	}

	return &packageSpec
}

func DeclarePackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	cmd.AddCommand(
		PackageInitCommand(),
		PackageUpCommand(),
		PackageDownCommand(),
		PackageRemoveCommand(),
	)

	return cmd
}
