package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
)

var ErrEmptyContainersObject = errors.New("empty supplied/returned container object")

func RemoveStaleInstantContainer(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Error(ctx, err)
	}

	for _, _container := range containers {
		for _, name := range _container.Names {
			if name == "/instant-openhie" {
				if _container.State == "running" {
					err = cli.ContainerStop(ctx, _container.ID, nil)
					if err != nil {
						log.Error(ctx, err)
					}
				}
				err = cli.ContainerRemove(ctx, _container.ID, types.ContainerRemoveOptions{})
				if err != nil {
					log.Error(ctx, err)
				}

				break
			}
		}
	}
}

func RemoveStaleInstantVolume(cli *client.Client, ctx context.Context) {
	volumes, err := cli.VolumeList(ctx, filters.Args{})
	if err != nil {
		log.Error(ctx, err)
	}

	for _, volume := range volumes.Volumes {
		if volume.Name == "instant" {
			err = cli.VolumeRemove(ctx, volume.Name, false)
			if err != nil {
				log.Error(ctx, err)
			}

			break
		}
	}
}

// This code attempts to combat old/dead containers lying around and being selected instead of the new container
func latestContainer(containers []types.Container, allowAllFails bool) (types.Container, error) {
	if len(containers) == 0 {
		return types.Container{}, errors.Wrap(ErrEmptyContainersObject, "")
	}

	var latestContainer types.Container
	for _, container := range containers {
		if container.Created > latestContainer.Created {
			latestContainer = container
		}
	}

	return latestContainer, nil
}

func ListContainerByName(containerName string) (types.Container, error) {
	client, err := NewDockerClient()
	if err != nil {
		return types.Container{}, err
	}

	filtersPair := filters.KeyValuePair{
		Key:   "name",
		Value: containerName,
	}

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filtersPair),
		All:     true,
	})
	if err != nil {
		return types.Container{}, errors.Wrap(err, "")
	}

	return latestContainer(containers, false)
}
