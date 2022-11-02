package pkg

import (
	viperUtil "cli/cmd/util"
	"cli/core"
	"cli/util"
	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidConfigFileSyntax = errors.New("invalid config file syntax, refer to https://github.com/openhie/package-starter-kit/blob/main/README.md, for information on valid config file syntax")
)

func setPackageActionFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringSliceP("name", "n", nil, "The name(s) of the package(s)")
	flags.BoolP("dev", "d", false, "For development related functionality (Passes `dev` as the second argument to your swarm file)")
	flags.BoolP("only", "o", false, "Ignore package dependencies")
	flags.StringP("profile", "p", "", "The profile name to load parameters from (defined in config.yml)")
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

	// TODO: don't panic on flag conflicts, the command line should override the profiles
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

func parseAndPrepareLaunch(cmd *cobra.Command) (*core.PackageSpec, *core.Config, error) {
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

	err = prepareEnvironment(*config)
	if err != nil {
		return nil, nil, err
	}

	return packageSpec, config, nil
}

func prepareEnvironment(config core.Config) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "")
	}

	core.RemoveStaleInstantContainer(cli, ctx)
	core.RemoveStaleInstantVolume(cli, ctx)

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return errors.Wrap(err, "")
	}

	if !hasImage(config.Image, images) {
		reader, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer reader.Close()

		// This io.Copy helps to wait for the image to finish downloading
		_, err = io.Copy(ioutil.Discard, reader)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func hasImage(imageName string, images []types.ImageSummary) bool {
	for _, image := range images {
		if util.SliceContains(image.RepoTags, imageName) {
			return true
		}
	}

	return false
}
