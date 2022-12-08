package util

import (
	"archive/tar"
	"archive/zip"
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
	"github.com/docker/docker/pkg/archive"
	"github.com/go-git/go-git/v5"
	"github.com/luno/jettison/errors"
)

var ErrEmptyContainersObject = errors.New("empty supplied/returned container object")

func CloneRepo(url, dest string) error {
	cloneOptions := &git.CloneOptions{
		URL: url,
	}

	_, err := git.PlainClone(dest, false, cloneOptions)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func CopyCredsToInstantContainer() (err error) {
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

	client, err := NewDockerClient()
	if err != nil {
		return err
	}

	instantContainer, err := listContainerByName("instant-openhie")
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

func listContainerByName(containerName string) (types.Container, error) {
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
	} else {
		clientOpts = append(clientOpts, client.FromEnv)
	}

	clientOpts = append(clientOpts, client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(clientOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return cli, nil
}

func unzipFile(f *zip.File, destination string) error {
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return errors.Wrap(errors.New("invalid file path: "+filePath), "")
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return errors.Wrap(err, "")
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return errors.Wrap(err, "")
	}

	// Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func UnzipSource(source, destination string) error {
	// Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer reader.Close()

	// Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return errors.Wrap(err, "")
	}

	// Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func UntarSource(source, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "")
	}

	tr := tar.NewReader(file)
	if err != nil {
		return errors.Wrap(err, "")
	}

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return errors.Wrap(err, "")

		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(destination, 0755)
			if err != nil {
				return errors.Wrap(err, "")
			}

		case tar.TypeReg:
			dir, _ := filepath.Split(destination)
			if dir != "" {
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return errors.Wrap(err, "")
				}
			}

			f, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrap(err, "")
			}

			if _, err := io.Copy(f, tr); err != nil {
				return errors.Wrap(err, "")
			}

			f.Close()
		}
	}
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

func SliceContains[Type comparable](slice []Type, element Type) bool {
	for _, s := range slice {
		if element == s {
			return true
		}
	}

	return false
}
