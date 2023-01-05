package generate

import (
	"embed"
	"html/template"
	"os"
	"path"

	"cli/core"

	"github.com/luno/jettison/errors"
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
