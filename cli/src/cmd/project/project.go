package project

import (
	viperUtil "cli/cmd/util"
	"cli/core"
	"context"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func getConfigFromParams(cmd *cobra.Command) (*core.Config, error) {
	var config core.Config
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	configViper, err := viperUtil.SetConfigViper(configFile)
	if err != nil {
		return nil, err
	}

	err = configViper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &config, nil
}

func getPackageSpecFromParams(cmd *cobra.Command) (*core.PackageSpec, error) {
	packageSpec := core.PackageSpec{}

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

	packageSpec = core.PackageSpec{
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

func getProjectAction(cmd *cobra.Command) (string, error) {
	action := ""
	switch true {
	case cmd.Flag("init").Changed:
		action = "init"
	case cmd.Flag("up").Changed:
		action = "up"
	case cmd.Flag("down").Changed:
		action = "down"
	case cmd.Flag("remove").Changed:
		action = "destroy"
	default:
		return action, errors.Wrap(errors.New("invalid action entered"), "")
	}

	return action, nil
}

func DeclareProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project level commands",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("init").Changed || cmd.Flag("up").Changed || cmd.Flag("down").Changed || cmd.Flag("remove").Changed {
				ctx := context.Background()

				packageSpec, err := getPackageSpecFromParams(cmd)
				if err != nil {
					log.Error(ctx, err)
					panic(err)
				}

				action, err := getProjectAction(cmd)
				if err != nil {
					log.Error(ctx, err)
					panic(err)
				}
				packageSpec.DeployCommand = action

				config, err := getConfigFromParams(cmd)
				if err != nil {
					log.Error(ctx, err)
					panic(err)
				}
				packageSpec.Packages = config.Packages
				packageSpec.CustomPackages = config.CustomPackages

				packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
				if err != nil {
					log.Error(ctx, err)
					panic(err)
				}

				err = core.LaunchDeploymentContainer(*packageSpec, *config)
				if err != nil {
					log.Error(ctx, err)
					panic(err)
				}
			}
		},
	}
	flags := cmd.Flags()

	flags.BoolP("init", "i", false, "Initialize all packages in the project")
	flags.BoolP("up", "u", false, "Up all packages in the project")
	flags.BoolP("down", "d", false, "Down all packages in the project")
	flags.BoolP("remove", "r", false, "Remove all packages in the project")

	cmd.MarkFlagsMutuallyExclusive("init", "up", "down", "remove")

	cmd.AddCommand(
		ProjectGenerateCommand(),
	)

	return cmd
}
