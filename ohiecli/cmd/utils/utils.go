package utils

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

var ErrEmptyContainersObject = errors.New("empty supplied/returned container object")

func GetHelpText(interactive bool, options string) string {
	if interactive {
		switch options {
		case "Deploy Commands":
			return `Commands:
				init/up/destroy/down	the deploy command you want to run (brief description below)
					deploy commands:
						init:			for initializing a service
						up:				for starting up a service that has been shut down or updating a service
						destroy:	for destroying a service
						down:			for bringing down a running service
			`
		case "Custom Options":
			return `Commands:
				Choose deploy action - for choosing the deploy action

				Choose target launcher - for choosing the deploy target. Can be a docker swarm, kubernetes or docker (Project specific)

				Specify deploy packages - for choosing the packages you want to use (core and custom packages)

				Specify environment variable file location - for specifying the file path to an environment variables file

				Specify environment variables - for specifying environment variables

				Specify custom package locations - for specifying the location or url to the custom packages you want to operate on

				Toggle only flag - for specifying the only flag, which specifies that actions are to be taken on a single package and not on its dependencies

				Specify Image Version - for specifying the version of the instant or platform image to use. Default is latest

				Toggle dev mode - for enabling the development mode in which the service ports are exposed

				Execute with current options - this executes the options that have been specified

				View current options set - for viewing the options that have been specified

				Reset to default options - for resetting to default options
			`
		default:
			return `Commands:
				Use Docker on your PC - this is for deploying packages to either docker or docker swarm

				Use a kubernetes Cluster - this is for deploying packages to a kubernetes cluster

				Install FHIR package - this is for installing FHIR IGs hosted remotely
			`
		}
	} else {
		return `Commands: 
		help 		this menu

		init/up/destroy/down	the deploy command you want to run (brief description below)
					deploy commands:
						init:	 for initializing a service
						up:	 for starting up a service that has been shut down or updating a service
						destroy: for destroying a service
						down:	 for bringing down a running service
					custom flags:
						--only, -o:							used to specify a single service for services that have dependencies. For cases where one wants to shut down or destroy a service without affecting its dependencies
						-t:											specifies the target to deploy to. Options are docker, swarm (docker swarm) and k8s (kubernetes) - project dependant.
						--custom-package, -c:		specifies path or url to a custom package. Git ssh urls are supported
						--dev:									specifies the development mode in which all service ports are exposed
						-e:											for specifying an environment variable
						--env-file: 						for specifying the path to an environment variables file
						--image-version:			the version of the project used for the deploy. Defaults to 'latest'
						-*, --*:								unrecognised flags are passed through uninterpreted
					usage:
						<deploy command> <custom flags> <package ids>
					examples:
						{your_binary_file} init -t=swarm --dev -e="NODE_ENV=prod" --env-file="../env.dev" -c="../customPackage1" -c="<git@github.com/customPackage2>"  interoperability-layer-openhim customPackage1_id customPackage2_id
						{your_binary_file} down -t=docker --only elastic_analytics

		install		install fhir npm package on fhir server
					usage: install <ig_url> <fhir_server> <authtype> <user/token> <pass>

					examples:
					install https://intrahealth.github.io/simple-hiv-ig/ http://hapi.fhir.org/baseR4 none
					install <ig_url> <fhir_server> basic smith stuff
					install <ig_url> <fhir_server> token "123"
					install <ig_url> <fhir_server> custom test
		`
	}
}

func SliceContains(slice []string, element string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}

func SliceContainsSubstr(slice []string, element string) bool {
	for _, s := range slice {
		if strings.Contains(element, s) {
			return true
		}
	}
	return false
}

func TarSource(path string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	ok := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "")
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.Wrap(err, "")
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, path, "", -1), string(filepath.Separator))
		err = tw.WriteHeader(header)
		if err != nil {
			return errors.Wrap(err, "")
		}

		f, err := os.Open(file)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if fi.IsDir() {
			return nil
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return errors.Wrap(err, "")
		}

		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "")
		}

		return nil
	})

	if ok != nil {
		return nil, ok
	}

	err := tw.Close()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return bufio.NewReader(&buf), nil
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

func NewDockerClient() (*client.Client, error) {
	var clientOpts []client.Opt

	host := os.Getenv("DOCKER_HOST")
	if host != "" {
		helper, err := connhelper.GetConnectionHelper(host)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}

		clientOpts = append(clientOpts,
			client.WithHTTPClient(httpClient),
			client.WithHost(helper.Host),
			client.WithDialContext(helper.Dialer),
		)
	}

	clientOpts = append(clientOpts, client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(clientOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return cli, nil
}
