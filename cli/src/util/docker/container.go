package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
)

var ErrEmptyContainersObject = errors.New("empty supplied/returned container object")

func RemoveStaleInstantContainer(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Error(ctx, err)
	}

	for _, _container := range containers {
		for _, name := range _container.Names {
			if name == "/instant-openhie" {
				if _container.State == "running" {
					err = cli.ContainerStop(ctx, _container.ID, container.StopOptions{})
					if err != nil {
						log.Error(ctx, err)
					}
				}
				err = cli.ContainerRemove(ctx, _container.ID, container.RemoveOptions{})
				if err != nil {
					log.Error(ctx, err)
				}

				break
			}
		}
	}
}

func RemoveStaleInstantVolume(cli *client.Client, ctx context.Context) {
	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
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
func latestContainer(containers []container.Summary, allowAllFails bool) (container.Summary, error) {
	if len(containers) == 0 {
		return container.Summary{}, errors.Wrap(ErrEmptyContainersObject, "")
	}

	var latestContainer container.Summary
	for _, container := range containers {
		if container.Created > latestContainer.Created {
			latestContainer = container
		}
	}

	return latestContainer, nil
}

func ListContainerByName(containerName string) (container.Summary, error) {
	client, err := NewDockerClient()
	if err != nil {
		return container.Summary{}, err
	}

	filtersPair := filters.KeyValuePair{
		Key:   "name",
		Value: containerName,
	}

	containers, err := client.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(filtersPair),
		All:     true,
	})
	if err != nil {
		return container.Summary{}, errors.Wrap(err, "")
	}

	return latestContainer(containers, false)
}
