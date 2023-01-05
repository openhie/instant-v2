package prompt

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/luno/jettison/errors"
	"github.com/manifoldco/promptui"
)

func GeneratePackagePrompt() (*generatePackagePromptResponse, error) {
	/*
		What is the id of your package:
		What is the name of your package:
		Provide a description of your package:
		Which stack does your package belong to:
		Do you want to include a dev compose file:
	*/

	promptId := promptui.Prompt{
		Label:   "What is the id of your package",
		Default: "my-package",
	}
	id, err := promptId.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptName := promptui.Prompt{
		Label:   "What is the name of your package",
		Default: "My Package",
	}
	name, err := promptName.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this project",
		Default: strcase.ToKebab(fmt.Sprintf("organisation/%v", name)),
	}
	image, err := promptImage.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptDescription := promptui.Prompt{
		Label:   "Provide a description of your package",
		Default: "A package to be used with the platform",
	}
	description, err := promptDescription.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptStack := promptui.Prompt{
		Label:   "Which stack does your package belong to",
		Default: "instant",
	}
	stack, err := promptStack.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptType := promptui.Select{
		Label: "What type best suites your package",
		Items: []string{"infrastructure", "use-case"},
	}
	index, packageType, err := promptType.Run()
	if err != nil || index == -1 {
		return nil, errors.Wrap(err, "")
	}

	promptDev := promptui.Select{
		Label: "Do you want to include a dev compose file",
		Items: []string{"Yes", "No"},
	}
	index, Dev, err := promptDev.Run()
	if err != nil || index == -1 {
		return nil, errors.Wrap(err, "")
	}

	var (
		includeDevFile bool
		targetPort     string
		publishedPort  string
	)
	if Dev == "Yes" {
		includeDevFile = true

		promptTargetPort := promptui.Prompt{
			Label:   "Which port would you like to target in dev mode?",
			Default: "8080",
		}
		targetPort, err = promptTargetPort.Run()
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		promptPublishedPort := promptui.Prompt{
			Label:   "Which port would you like published in dev mode?",
			Default: "8081",
		}
		publishedPort, err = promptPublishedPort.Run()
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	promptResponse := &generatePackagePromptResponse{
		Id:             id,
		Name:           name,
		Image:          image,
		Description:    description,
		Stack:          stack,
		Type:           packageType,
		IncludeDevFile: includeDevFile,
		TargetPort:     targetPort,
		PublishedPort:  publishedPort,
	}

	return promptResponse, nil
}
