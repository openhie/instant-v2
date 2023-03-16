package parse

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"cli/core"
	"cli/core/state"
	"cli/util/docker"
	"cli/util/slice"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

func ParseAndPrepareLaunch(cmd *cobra.Command) (*core.PackageSpec, *core.Config, error) {
	if cmd.Flags().Changed("env-file") {
		envFiles, err := cmd.Flags().GetStringSlice("env-file")
		if err != nil {
			return nil, nil, errors.Wrap(err, "")
		}

		state.EnvFiles = nil

		for _, envFile := range envFiles {
			if !filepath.IsAbs(envFile) {
				absFilePath, err := filepath.Abs(envFile)
				if err != nil {
					return nil, nil, errors.Wrap(err, "")
				}

				state.EnvFiles = append(state.EnvFiles, absFilePath)
				continue
			}

			state.EnvFiles = append(state.EnvFiles, envFile)
		}
	}

	config, err := GetConfigFromParams(cmd)
	if err != nil {
		return nil, nil, err
	}

	packageSpec, err := getPackageSpecFromParams(cmd, config)
	if err != nil {
		return nil, nil, err
	}

	packageSpec, err = getPackageSpecFromProfile(cmd, *config, *packageSpec)
	if err != nil {
		return nil, nil, err
	}

	packageSpec, err = filterEnvVars(cmd, packageSpec)
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

	cli, err := docker.NewDockerClient()
	if err != nil {
		return errors.Wrap(err, "")
	}

	docker.RemoveStaleInstantContainer(cli, ctx)
	docker.RemoveStaleInstantVolume(cli, ctx)

	hasImage, err := hasImage(cli, config.Image)
	if err != nil {
		return err
	}

	if !hasImage {
		fmt.Println("> Image", config.Image, "can't be found locally .. Pulling from docker")
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

func hasImage(dockerCli *client.Client, imageName string) (bool, error) {
	images, err := dockerCli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, errors.Wrap(err, "")
	}

	for _, image := range images {
		if slice.SliceContains(image.RepoTags, imageName) {
			return true, nil
		}
	}

	return false, nil
}

// filterEnvVars checks for env vars that could have been set be prior functions (like in --profile),
// and ensures that env vars parsed are taken in order of precedence of --env-var >> --env-file >> profile env files
func filterEnvVars(cmd *cobra.Command, pSpec *core.PackageSpec) (*core.PackageSpec, error) {
	var paramsEnvVars []string
	if cmd.Flags().Changed("env-var") {
		envVars, err := cmd.Flags().GetStringSlice("env-var")
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		paramsEnvVars = envVars
	}

	envVarsMap := make(map[string]string)
	if len(paramsEnvVars) > 0 {
		envVarsMap = slice.AppendUniqueToMapFromSlice(envVarsMap, paramsEnvVars)
	}

	var paramsEnvFileEnvVars []string
	if cmd.Flags().Changed("env-file") {
		envFiles, err := cmd.Flags().GetStringSlice("env-file")
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		envViper, err := state.GetEnvironmentVariableViper(envFiles)
		if err != nil {
			return nil, err
		}
		paramsEnvFileEnvVars = state.GetEnvVariableString(envViper)
	}
	if len(paramsEnvFileEnvVars) > 0 {
		envVarsMap = slice.AppendUniqueToMapFromSlice(envVarsMap, paramsEnvFileEnvVars)
	}
	if len(pSpec.EnvironmentVariables) > 0 {
		envVarsMap = slice.AppendUniqueToMapFromSlice(envVarsMap, pSpec.EnvironmentVariables)
	}

	pSpec.EnvironmentVariables = []string{}
	for k, v := range envVarsMap {
		pSpec.EnvironmentVariables = append(pSpec.EnvironmentVariables, k+v)
	}

	return pSpec, nil
}
