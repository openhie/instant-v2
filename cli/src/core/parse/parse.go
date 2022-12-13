package parse

import (
	"context"
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

func parseAndPrepareLaunch(cmd *cobra.Command) (*core.PackageSpec, *core.Config, error) {
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

	packageSpec, err := GetPackageSpecFromParams(cmd, config)
	if err != nil {
		return nil, nil, err
	}

	packageSpec, err = GetPackageSpecFromProfile(cmd, *config, *packageSpec)
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
