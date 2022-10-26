package pkg

import (
	viperUtil "cli/cmd/util"
	"cli/core"
	"context"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrConflictingDevFlag      = errors.New("conflicting command-line and profile flag: --dev")
	ErrConflictingOnlyFlag     = errors.New("conflicting command-line and profile flag: --only")
	ErrInvalidConfigFileSyntax = errors.New("invalid config file syntax, refer to https://github.com/openhie/package-starter-kit/blob/main/README.md, for information on valid config file syntax")
)

func setPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.Bool("dev", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.Bool("only", false, "Ignore package dependencies")
	flags.String("profile", "", "The profile name to load parameters from (defined in config.yml)")
	flags.StringSliceP("custom-path", "c", nil, "Path(s) to custom package(s)")

	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getConfigFromParams(cmd)
		if err != nil {
			log.Error(context.Background(), err)
		}

		return config.Packages, cobra.ShellCompDirectiveDefault
	})
	cmd.RegisterFlagCompletionFunc("profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := getConfigFromParams(cmd)
		if err != nil {
			log.Error(context.Background(), err)
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
			log.Error(context.Background(), err)
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
		return nil, errors.Wrap(err, "")
	}

	configViper, err := viperUtil.GetConfigViper(configFile)
	if err != nil {
		return nil, err
	}

	populatedConfig, err := unmarshalConfig(config, configViper)
	if err != nil {
		return nil, err
	}

	appendTag(populatedConfig)

	return populatedConfig, nil
}

func appendTag(config *core.Config) {
	splitStrings := strings.Split(config.Image, ":")

	if len(splitStrings) == 1 {
		config.Image += ":latest"
	}
}

func unmarshalConfig(config core.Config, configViper *viper.Viper) (*core.Config, error) {
	err := configViper.Unmarshal(&config)
	if err != nil && strings.Contains(err.Error(), "expected type") {
		return nil, errors.Wrap(ErrInvalidConfigFileSyntax, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &config, nil
}

func getCustomPackages(config *core.Config, customPackagePaths []string) []core.CustomPackage {
	var customPackages []core.CustomPackage
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
				Path: customPackagePath,
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
		return nil, errors.Wrap(err, "")
	}
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	isDev, err := cmd.Flags().GetBool("dev")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	isOnly, err := cmd.Flags().GetBool("only")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	envFiles, err := cmd.Flags().GetStringSlice("env-file")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	envViper, err := viperUtil.GetEnvironmentVariableViper(envFiles)
	if err != nil {
		return nil, err
	}
	envVariables := viperUtil.GetEnvVariableString(envViper)

	customPackages := getCustomPackages(config, customPackagePaths)

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
		return nil, errors.Wrap(err, "")
	}

	for _, p := range config.Profiles {
		if p.Name == profileName {
			profile = p
			break
		}
	}

	if !cmd.Flags().Changed("dev") && profile.Dev {
		packageSpec.IsDev = profile.Dev
	} else if cmd.Flags().Changed("dev") && profileName != "" {
		val, err := cmd.Flags().GetBool("dev")
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		if val != profile.Dev {
			return nil, errors.Wrap(ErrConflictingDevFlag, "")
		}
	}

	if !cmd.Flags().Changed("only") && profile.Only {
		packageSpec.IsOnly = profile.Only
	} else if cmd.Flags().Changed("only") && profileName != "" {
		val, err := cmd.Flags().GetBool("only")
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		if val != profile.Only {
			return nil, errors.Wrap(ErrConflictingOnlyFlag, "")
		}
	}

	if len(profile.Packages) > 0 {
		packageSpec.Packages = append(profile.Packages, packageSpec.Packages...)
	}
	if len(profile.EnvFiles) > 0 {
		envViper, err := viperUtil.GetEnvironmentVariableViper(profile.EnvFiles)
		if err != nil {
			return nil, err
		}
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

func packageActionHook(cmd *cobra.Command, args []string) (*core.PackageSpec, *core.Config, error) {
	config, err := getConfigFromParams(cmd)
	if err != nil {
	return nil, nil, err
	}

	packageSpec, err := getPackageSpecFromParams(cmd, config)
	if err != nil {
	return nil, nil, err
	}

	packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
	if err != nil {
	return nil, nil, err
	}

	err = validate(cmd, config)
	if err != nil {
	return nil, nil, err
	}
	
	for _, pack := range packageSpec.Packages {
		for _, customPack := range config.CustomPackages {
			if pack == customPack.Id {
				packageSpec.CustomPackages = append(packageSpec.CustomPackages, customPack)
			}
		}
	}

	return packageSpec, config, nil
}
