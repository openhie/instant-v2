package prompt

import (
	"cli/core"

	"github.com/luno/jettison/errors"
	"github.com/manifoldco/promptui"
)

func GeneratePackagePrompt() (core.GeneratePackageSpec, error) {
	/*
		What is the id of your package:
		What is the name of your package:
		What docker image would you like to use with this project:
		Provide a description of your package:
		What type best suites your package:
		Do you want to include a dev compose file:
		Which port would you like to target on the container in dev mode:
		Which port would you like published on the host in dev mode:
	*/

	promptId := promptui.Prompt{
		Label:   "What is the id of your package",
		Default: "my-package",
	}
	id, err := promptId.Run()
	if err != nil {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	promptName := promptui.Prompt{
		Label:   "What is the name of your package",
		Default: "My Package",
	}
	name, err := promptName.Run()
	if err != nil {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	promptImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this package",
		Default: "nginx",
	}
	image, err := promptImage.Run()
	if err != nil {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	promptDescription := promptui.Prompt{
		Label:   "Provide a description of your package",
		Default: "A package to be used with the platform",
	}
	description, err := promptDescription.Run()
	if err != nil {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	promptType := promptui.Select{
		Label: "What type best suites your package",
		Items: []string{"infrastructure", "use-case"},
	}
	index, packageType, err := promptType.Run()
	if err != nil || index == -1 {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	promptDev := promptui.Select{
		Label: "Do you want to include a dev compose file",
		Items: []string{"Yes", "No"},
	}
	index, Dev, err := promptDev.Run()
	if err != nil || index == -1 {
		return core.GeneratePackageSpec{}, errors.Wrap(err, "")
	}

	var (
		includeDevFile bool
		targetPort     string
		publishedPort  string
	)
	if Dev == "Yes" {
		includeDevFile = true

		promptTargetPort := promptui.Prompt{
			Label:   "Which port would you like to target on the container in dev mode?",
			Default: "80",
		}
		targetPort, err = promptTargetPort.Run()
		if err != nil {
			return core.GeneratePackageSpec{}, errors.Wrap(err, "")
		}

		promptPublishedPort := promptui.Prompt{
			Label:   "Which port would you like published on the host in dev mode?",
			Default: "8080",
		}
		publishedPort, err = promptPublishedPort.Run()
		if err != nil {
			return core.GeneratePackageSpec{}, errors.Wrap(err, "")
		}
	}

	promptResponse := core.GeneratePackageSpec{
		Id:             id,
		Name:           name,
		Image:          image,
		Description:    description,
		Type:           packageType,
		IncludeDevFile: includeDevFile,
		TargetPort:     targetPort,
		PublishedPort:  publishedPort,
	}

	return promptResponse, nil
}
