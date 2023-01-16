package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/log"
)

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
