package core

import (
	"context"
	"embed"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"cli/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

var TemplateFs embed.FS

type Profile struct {
	Name     string   `yaml:"name"`
	EnvFiles []string `yaml:"envFiles"`
	Dev      bool     `yaml:"dev"`
	Only     bool     `yaml:"only"`
	Packages []string `yaml:"packages"`
}

type CustomPackage struct {
	Id          string `yaml:"id"`
	Path        string `yaml:"path"`
	SshKey      string `yaml:"sshKey"`
	SshPassword string `yaml:"sshPassword"`
}

type Config struct {
	Image          string          `yaml:"image"`
	LogPath        string          `yaml:"logPath"`
	Packages       []string        `yaml:"packages"`
	CustomPackages []CustomPackage `yaml:"customPackages"`
	Profiles       []Profile       `yaml:"profiles"`
	ProjectName    string          `yaml:"projectName"`
	PlatformImage  string          `yaml:"platformImage"`
}

type PackageSpec struct {
	EnvironmentVariables []string
	DeployCommand        string
	Packages             []string
	IsDev                bool
	IsOnly               bool
	CustomPackages       []CustomPackage
	ImageVersion         string
	TargetLauncher       string
}

func getCustomPackageName(customPackage CustomPackage) string {
	if customPackage.Id != "" {
		return customPackage.Id
	}
	return strings.TrimSuffix(path.Base(path.Clean(customPackage.Path)), path.Ext(customPackage.Path))
}

func mountCustomPackage(customPackage CustomPackage, cli *client.Client, ctx context.Context, instantContainerId string) error {
	gitRegex := regexp.MustCompile(`\.git`)
	httpRegex := regexp.MustCompile("http")
	zipRegex := regexp.MustCompile(`\.zip`)
	tarRegex := regexp.MustCompile(`\.tar`)

	const CUSTOM_PACKAGE_LOCAL_PATH = "/tmp/custom-package/"
	customPackageName := getCustomPackageName(customPackage)
	customPackageTmpLocation := path.Join(CUSTOM_PACKAGE_LOCAL_PATH, customPackageName)
	err := os.RemoveAll(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return err
	}
	err = os.MkdirAll(customPackageTmpLocation, os.ModePerm)
	if err != nil {
		return err
	}

	if gitRegex.MatchString(customPackage.Path) {

		err = util.CloneRepo(customPackage.Path, customPackageTmpLocation, customPackage.SshKey, customPackage.SshPassword)
		if err != nil {
			return err
		}
	} else if httpRegex.MatchString(customPackage.Path) {
		resp, err := http.Get(customPackage.Path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.Wrapf(err, "Error in downloading custom package - HTTP status code: %v", strconv.Itoa(resp.StatusCode))
		}

		if zipRegex.MatchString(customPackage.Path) {
			tmpZip, err := os.CreateTemp("", "tmp-*.zip")
			if err != nil {
				return err
			}

			_, err = io.Copy(tmpZip, resp.Body)
			if err != nil {
				return err
			}
			err = util.UnzipSource(tmpZip.Name(), customPackageTmpLocation)
			if err != nil {
				return err
			}
			err = os.Remove(tmpZip.Name())
			if err != nil {
				return err
			}
		} else if tarRegex.MatchString(customPackage.Path) {
			tmpTar, err := os.CreateTemp("", "tmp-*.tar")
			if err != nil {
				return err
			}

			_, err = io.Copy(tmpTar, resp.Body)
			if err != nil {
				return err
			}

			err = util.UntarSource(tmpTar.Name(), customPackageTmpLocation)
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
			return err
		}
	}

	customPackageReader, err := util.TarSource(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return err
	}
	err = cli.CopyToContainer(ctx, instantContainerId, "/instant/", customPackageReader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	err = os.RemoveAll(CUSTOM_PACKAGE_LOCAL_PATH)
	if err != nil {
		return err
	}

	return nil
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

// Attaches a container's STDOUT until that container has been removed
func attachUntilRemoved(cli *client.Client, ctx context.Context, instantContainerId string) error {
	attachResponse, err := cli.ContainerAttach(ctx, instantContainerId, types.ContainerAttachOptions{Stdout: true, Stream: true, Logs: true, Stderr: true})
	if err != nil {
		return err
	}
	defer attachResponse.Close()

	go func() {
		_, err = stdcopy.StdCopy(os.Stdout, os.Stdout, attachResponse.Reader)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()

	successC, errC := cli.ContainerWait(ctx, instantContainerId, "removed")
	select {
	case <-successC:
		return nil
	case err := <-errC:
		return err
	}
}

func getInstantCommand(packageSpec PackageSpec) []string {
	instantCommand := []string{packageSpec.DeployCommand, "-t", "swarm"}

	if packageSpec.IsDev {
		instantCommand = append(instantCommand, "--dev")
	}

	if packageSpec.IsOnly {
		instantCommand = append(instantCommand, "--only")
	}

	instantCommand = append(instantCommand, packageSpec.Packages...)

	for _, customPackage := range packageSpec.CustomPackages {
		instantCommand = append(instantCommand, getCustomPackageName(customPackage))
	}

	return instantCommand
}

func LaunchPackage(packageSpec PackageSpec, config Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	removeStaleInstantContainer(cli, ctx)
	removeStaleInstantVolume(cli, ctx)

	reader, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
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

	instantCommand := getInstantCommand(packageSpec)

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
		return err
	}

	for _, customPackage := range packageSpec.CustomPackages {
		err = mountCustomPackage(customPackage, cli, ctx, instantContainer.ID)
		if err != nil {
			return err
		}
	}

	defer removeStaleInstantVolume(cli, ctx)

	err = cli.ContainerStart(ctx, instantContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	err = attachUntilRemoved(cli, ctx, instantContainer.ID)
	if err != nil {
		return err
	}

	return nil
}

type GeneratePackageSpec struct {
	Id             string
	Name           string
	Stack          string
	Description    string
	Type           string
	IncludeDevFile bool
}

func createFileFromTemplate(source, destination string, generatePackageSpec GeneratePackageSpec) error {
	destination = path.Join(destination, source)
	templatePath := path.Join("template/package/", source)

	packageTemplate, err := template.New("package").ParseFS(TemplateFs, templatePath)
	if err != nil {
		return err
	}
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	err = packageTemplate.ExecuteTemplate(file, source, generatePackageSpec)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePackage(destination string, generatePackageSpec GeneratePackageSpec) error {

	createFileFromTemplate("swarm.sh", destination, generatePackageSpec)
	createFileFromTemplate("package-metadata.json", destination, generatePackageSpec)
	createFileFromTemplate("docker-compose.yml", destination, generatePackageSpec)

	if generatePackageSpec.IncludeDevFile {
		createFileFromTemplate("docker-compose.dev.yml", destination, generatePackageSpec)
	}

	return nil
}
