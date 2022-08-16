package core

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
)

type Package struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type Config struct {
	Image    string    `yaml:"image"`
	LogPath  string    `yaml:"logPath"`
	Packages []Package `yaml:"packages"`
}

type PackageSpec struct {
	EnvironmentVariables []string
	DeployCommand        string
	Packages             []string
	IsDev                bool
	IsOnly               bool
	CustomPackagePaths   []string
	ImageVersion         string
	TargetLauncher       string
}

func removeStaleInstantContainer(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	util.LogError(err)

	for _, _container := range containers {
		for _, name := range _container.Names {
			if name == "/instant-openhie" {
				if _container.State == "running" {
					err = cli.ContainerStop(ctx, _container.ID, nil)
					util.PanicError(err)
				}
				err = cli.ContainerRemove(ctx, _container.ID, types.ContainerRemoveOptions{})
				util.LogError(err)
				break
			}
		}
	}
}

func removeStaleInstantVolume(cli *client.Client, ctx context.Context) {
	volumes, err := cli.VolumeList(ctx, filters.Args{})
	util.LogError(err)

	for _, volume := range volumes.Volumes {
		if volume.Name == "instant" {
			err = cli.VolumeRemove(ctx, volume.Name, false)
			util.LogError(err)
			break
		}
	}
}

func attachStdoutToInstantOutput(cli *client.Client, ctx context.Context, instantContainerId string) {
	attachResponse, err := cli.ContainerAttach(ctx, instantContainerId, types.ContainerAttachOptions{Stdout: true, Stream: true, Logs: true, Stderr: true})
	util.PanicError(err)
	defer attachResponse.Close()
	os.Stdout.ReadFrom(attachResponse.Reader)
}

func LaunchPackage(packageSpec PackageSpec, config Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	util.PanicError(err)

	removeStaleInstantContainer(cli, ctx)
	removeStaleInstantVolume(cli, ctx)

	reader, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	util.PanicError(err)
	defer reader.Close()
	if os.Getenv("LOG") == "true" {
		io.Copy(os.Stdout, reader)
	}

	mounts := []mount.Mount{
		{
			Type:   mount.TypeVolume,
			Source: "instant",
			Target: "/instant",
		},
	}

	if config.LogPath != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: config.LogPath,
			Target: "/tmp/logs",
		})
	}

	endpointSettings := make(map[string]*network.EndpointSettings)
	endpointSettings["host"] = &network.EndpointSettings{
		EndpointID: "host",
	}

	instantCommand := []string{packageSpec.DeployCommand, "-t", "swarm"}
	if packageSpec.IsDev {
		instantCommand = append(instantCommand, "--dev")
	}
	if packageSpec.IsOnly {
		instantCommand = append(instantCommand, "--only")
	}
	instantCommand = append(instantCommand, packageSpec.Packages...)

	instantContainer, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        config.Image,
		Cmd:          instantCommand,
		AttachStderr: true,
		AttachStdout: true,
		Env:          packageSpec.EnvironmentVariables,
	}, &container.HostConfig{
		NetworkMode: "host",
		Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		Mounts:      mounts,
	}, &network.NetworkingConfig{EndpointsConfig: endpointSettings}, nil, "instant-openhie")
	util.PanicError(err)

	err = cli.ContainerStart(ctx, instantContainer.ID, types.ContainerStartOptions{})
	util.PanicError(err)

	attachStdoutToInstantOutput(cli, ctx, instantContainer.ID)

	successC, errC := cli.ContainerWait(ctx, instantContainer.ID, "exited")
	select {
	case <-successC:
		err = cli.ContainerRemove(ctx, instantContainer.ID, types.ContainerRemoveOptions{})
		util.LogError(err)
		removeStaleInstantVolume(cli, ctx)
	case <-errC:
		util.PanicError(err)
	}

}
