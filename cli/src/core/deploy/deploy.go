package deploy

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"

	"cli/core"
	"cli/core/parse"
	"cli/util/file"
	"cli/util/git"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	cp "github.com/otiai10/copy"
)

func MountCustomPackage(ctx context.Context, cli *client.Client, customPackage core.CustomPackage, instantContainerId string) error {
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

func LaunchDeploymentContainer(packageSpec *core.PackageSpec, config *core.Config) error {
	// TODO: code to launch instant container

	return nil
}
