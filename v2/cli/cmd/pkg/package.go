package pkg

import (
	"log"

	viperUtil "github.com/openhie/package-starter-kit/cli/v2/cli/cmd/util"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func setPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("dev", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.Bool("only", false, "Ignore package dependencies")
	flags.String("profile", "", "The profile name to load parameters from (defined in config.yml)")
	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")
	flags.String("ssh-key", "", "The path to the ssh key required for cloning a custom package")
	flags.String("ssh-password", "", "The password (or path to the file containing the password) required for authenticating the ssh-key when cloning a custom package")

	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getConfigFromParams(cmd)
		if err != nil {
			log.Print(err)
		}
		return config.Packages, cobra.ShellCompDirectiveDefault
	})
	cmd.RegisterFlagCompletionFunc("profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getConfigFromParams(cmd)
		if err != nil {
			log.Print(err)
		}
		var profileNames []string
		for _, p := range config.Profiles {
			profileNames = append(profileNames, p.Name)
		}
		return profileNames, cobra.ShellCompDirectiveDefault
	})
	cmd.RegisterFlagCompletionFunc("custom-path", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getConfigFromParams(cmd)
		if err != nil {
			log.Print(err)
		}
		var customPackages []string
		for _, c := range config.CustomPackages {
			customPackages = append(customPackages, c.Id)
		}
		return customPackages, cobra.ShellCompDirectiveDefault
	})
}

func getConfigFromParams(cmd *cobra.Command) (*core.Config, error) {
	var config core.Config
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}
	configViper := viperUtil.GetConfigViper(configFile)
	err = configViper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func getCustomPackages(config *core.Config, customPackagePaths []string, sshKey, sshPassword string) []core.CustomPackage {
	customPackages := make([]core.CustomPackage, len(customPackagePaths))

	for _, customPackagePath := range customPackagePaths {
		var customPackage core.CustomPackage

		for _, configCustomPackage := range config.CustomPackages {
			if customPackagePath == configCustomPackage.Id || customPackagePath == configCustomPackage.Path {
				customPackage = configCustomPackage
				break
			}
		}
		if customPackage.Id == "" {
			customPackage = core.CustomPackage{
				Path:        customPackagePath,
				SshKey:      sshKey,
				SshPassword: sshPassword,
			}
		}

		customPackages = append(customPackages, customPackage)
	}
	return customPackages
}

func getPackageSpecFromParams(cmd *cobra.Command, config *core.Config) (*core.PackageSpec, error) {
	packageSpec := core.PackageSpec{}

	packageNames, err := cmd.Flags().GetStringSlice("name")
	if err != nil {
		return nil, err
	}
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return nil, err
	}
	isDev, err := cmd.Flags().GetBool("dev")
	if err != nil {
		return nil, err
	}
	isOnly, err := cmd.Flags().GetBool("only")
	if err != nil {
		return nil, err
	}

	envFiles, err := cmd.Flags().GetStringSlice("env-file")
	if err != nil {
		return nil, err
	}
	envViper := viperUtil.GetEnvironmentVariableViper(envFiles)
	envVariables := viperUtil.GetEnvVariableString(envViper)

	sshKey, err := cmd.Flags().GetString("ssh-key")
	util.PanicError(err)

	sshPassword, err := cmd.Flags().GetString("ssh-password")
	util.PanicError(err)

	customPackages := getCustomPackages(config, customPackagePaths, sshKey, sshPassword)

	packageSpec = core.PackageSpec{
		Packages:             packageNames,
		CustomPackages:       customPackages,
		EnvironmentVariables: envVariables,
		IsDev:                isDev,
		IsOnly:               isOnly,
		DeployCommand:        cmd.Use,
	}

	return &packageSpec, nil
}

func loadInProfileParams(cmd *cobra.Command, config core.Config, packageSpec core.PackageSpec) (*core.PackageSpec, error) {
	profile := core.Profile{}
	profileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		return nil, err
	}
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

	return &packageSpec, nil
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
		PackageGenerateCommand(),
	)

	return cmd
}
