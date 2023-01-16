package generate

import (
	"embed"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"cli/core"

	"github.com/luno/jettison/errors"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed template/*
	templateFs       embed.FS
	ErrInvalidConfig = errors.New("invalid project config, required fields are Image, ProjectName, and PlatformImage")
)

func createFileFromTemplate(source, destination string, generatePackageSpec core.GeneratePackageSpec) error {
	destination = path.Join(destination, source)
	templatePath := path.Join("template", "package", source)

	packageTemplate, err := template.New("package").ParseFS(templateFs, templatePath)
	if err != nil {
		return errors.Wrap(err, "")
	}

	file, err := os.Create(destination)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = packageTemplate.ExecuteTemplate(file, source, generatePackageSpec)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func GeneratePackage(destination string, generatePackageSpec core.GeneratePackageSpec) error {
	err := createFileFromTemplate("swarm.sh", destination, generatePackageSpec)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = createFileFromTemplate("package-metadata.json", destination, generatePackageSpec)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = createFileFromTemplate("docker-compose.yml", destination, generatePackageSpec)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if generatePackageSpec.IncludeDevFile {
		err = createFileFromTemplate("docker-compose.dev.yml", destination, generatePackageSpec)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}

func GenerateConfigFile(config *core.Config) error {
	if config.Image == "" || config.ProjectName == "" || config.PlatformImage == "" {
		return errors.Wrap(ErrInvalidConfig, "")
	}

	firstFields := core.Config{
		ProjectName:   config.ProjectName,
		Image:         config.Image,
		PlatformImage: config.PlatformImage,
		LogPath:       config.LogPath,
	}

	data, err := yaml.Marshal(&firstFields)
	if err != nil {
		return errors.Wrap(err, "")
	}
	data = append(data, '\n')

	secondFields := core.Config{
		Packages: config.Packages,
	}

	d, err := yaml.Marshal(&secondFields)
	if err != nil {
		return errors.Wrap(err, "")
	}
	data = append(data, d...)
	data = append(data, '\n')

	thirdFields := core.Config{
		CustomPackages: config.CustomPackages,
	}

	d, err = yaml.Marshal(&thirdFields)
	if err != nil {
		return errors.Wrap(err, "")
	}
	data = append(data, d...)
	data = append(data, '\n')

	fourthFields := core.Config{
		Profiles: config.Profiles,
	}

	d, err = yaml.Marshal(&fourthFields)
	if err != nil {
		return errors.Wrap(err, "")
	}
	data = append(data, d...)

	err = ioutil.WriteFile("config.yaml", data, 0600)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
