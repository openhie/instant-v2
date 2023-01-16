package deploy

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cli/core"
	"cli/core/parse"
	"cli/util/docker"
	"cli/util/file"
	"cli/util/git"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/log"
	cp "github.com/otiai10/copy"
)

func mountCustomPackage(ctx context.Context, cli *client.Client, customPackage core.CustomPackage, instantContainerId string) error {
	gitRegex := regexp.MustCompile(`\.git`)
	httpRegex := regexp.MustCompile("http")
	zipRegex := regexp.MustCompile(`\.zip`)
	tarRegex := regexp.MustCompile(`\.tar`)

	const CUSTOM_PACKAGE_LOCAL_PATH = "/tmp/custom-package/"
	customPackageTmpLocation := path.Join(CUSTOM_PACKAGE_LOCAL_PATH, parse.GetCustomPackageName(customPackage))
	err := os.RemoveAll(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = os.MkdirAll(customPackageTmpLocation, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if gitRegex.MatchString(customPackage.Path) && !httpRegex.MatchString(customPackage.Path) {
		err = git.CloneRepo(customPackage.Path, customPackageTmpLocation)
		if err != nil {
			return err
		}

	} else if httpRegex.MatchString(customPackage.Path) {
		resp, err := http.Get(customPackage.Path)
		if err != nil {
			return errors.Wrap(err, "")
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.Wrap(err, "Error in downloading custom package - HTTP status code: "+strconv.Itoa(resp.StatusCode))
		}

		if zipRegex.MatchString(customPackage.Path) {
			tmpZip, err := os.CreateTemp("", "tmp-*.zip")
			if err != nil {
				return errors.Wrap(err, "")
			}

			_, err = io.Copy(tmpZip, resp.Body)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = file.UnzipSource(tmpZip.Name(), customPackageTmpLocation)
			if err != nil {
				return err
			}

			err = os.Remove(tmpZip.Name())
			if err != nil {
				return errors.Wrap(err, "")
			}

		} else if tarRegex.MatchString(customPackage.Path) {
			tmpTar, err := os.CreateTemp("", "tmp-*.tar")
			if err != nil {
				return errors.Wrap(err, "")
			}

			_, err = io.Copy(tmpTar, resp.Body)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = file.UntarSource(tmpTar.Name(), customPackageTmpLocation)
			if err != nil {
				return err
			}

			err = os.Remove(tmpTar.Name())
			if err != nil {
				return err
			}
		}
	} else {
		err := cp.Copy(customPackage.Path, customPackageTmpLocation)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	customPackageReader, err := file.TarSource(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return err
	}
	err = cli.CopyToContainer(ctx, instantContainerId, "/instant/", customPackageReader, types.CopyToContainerOptions{})
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = os.RemoveAll(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func copyCredsToInstantContainer() (err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "")
	}
	dockerCredsPath := filepath.Join(homeDir, ".docker", "config.json")

	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "")
	} else if os.IsNotExist(err) {
		return nil
	}

	client, err := docker.NewDockerClient()
	if err != nil {
		return err
	}

	instantContainer, err := docker.ListContainerByName("instant-openhie")
	if err != nil {
		return err
	}

	dstInfo := archive.CopyInfo{
		Path:   "/root/.docker/",
		Exists: true,
		IsDir:  true,
	}

	srcInfo, err := archive.CopyInfoSourcePath(dockerCredsPath, false)
	if err != nil {
		return errors.Wrap(err, "")
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer srcArchive.Close()

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer preparedArchive.Close()

	err = client.CopyToContainer(context.Background(), instantContainer.ID, dstDir, preparedArchive, types.CopyToContainerOptions{
		CopyUIDGID: true,
	})
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

// Attaches a container's STDOUT until that container has been removed
func attachUntilRemoved(cli client.ContainerAPIClient, ctx context.Context, instantContainerId string) error {
	attachResponse, err := cli.ContainerAttach(ctx, instantContainerId, types.ContainerAttachOptions{Stdout: true, Stream: true, Logs: true, Stderr: true})
	if err != nil {
		return err
	}
	defer attachResponse.Close()

	go func() {
		_, err = stdcopy.StdCopy(os.Stdout, os.Stdout, attachResponse.Reader)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Error(ctx, err)
			panic(err)
		}
	}()

	successC, errC := cli.ContainerWait(ctx, instantContainerId, "removed")
	select {
	case <-successC:
		return nil
	case err := <-errC:
		if strings.Contains(err.Error(), "No such container") {
			return nil
		}
		return errors.Wrap(err, "")
	}
}

func LaunchDeploymentContainer(packageSpec *core.PackageSpec, config *core.Config) error {
	ctx := context.Background()

	cli, err := docker.NewDockerClient()
	if err != nil {
		return errors.Wrap(err, "")
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

	instantCommand := parse.GetInstantCommand(*packageSpec)

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
		AutoRemove:  true,
	}, &network.NetworkingConfig{EndpointsConfig: endpointSettings}, nil, "instant-openhie")
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, customPackage := range packageSpec.CustomPackages {
		err = mountCustomPackage(ctx, cli, customPackage, instantContainer.ID)
		if err != nil {
			return err
		}
	}

	err = copyCredsToInstantContainer()
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, instantContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = attachUntilRemoved(cli, ctx, instantContainer.ID)
	if err != nil {
		return err
	}

	return nil
}
