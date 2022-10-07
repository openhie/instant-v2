package pkg

import (
	"context"
	"reflect"

	viperUtil "cli/cmd/util"
	"cli/core"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// TODO: This function MUST be unit-tested
func getConfigFromParams(cmd *cobra.Command) (*core.Config, error) {
	var config core.Config
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	configViper, err := viperUtil.GetConfigViper(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var decoderOptions viper.DecoderConfigOption = func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = func(k1, k2 reflect.Kind, i interface{}) (interface{}, error) {
			if k1 == reflect.Map {
				ip := i.(map[string]interface{})

				// TODO: implement better logic for this
				_, sshKeyExists := ip["sshkey"]
				if _, ok := ip["packages"]; !ok && !sshKeyExists {
					return ip["id"], nil
				}
			}

			return i, nil
		}
	}

	err = configViper.Unmarshal(&config, decoderOptions)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &config, nil
}

func getCustomPackages(config *core.Config, customPackagePaths []string, sshKey, sshPassword string) []core.CustomPackage {
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

	sshKey, err := cmd.Flags().GetString("ssh-key")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	sshPassword, err := cmd.Flags().GetString("ssh-password")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

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
	}
	if !cmd.Flags().Changed("only") && profile.Only {
		packageSpec.IsOnly = profile.Only
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
