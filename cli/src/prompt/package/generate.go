package prompt

import (
	"github.com/luno/jettison/errors"
	"github.com/manifoldco/promptui"
)

type GeneratePackagePromptResponse struct {
	Id             string
	Name           string
	Stack          string
	Description    string
	Type           string
	IncludeDevFile bool
}

func GeneratePackagePrompt() (*GeneratePackagePromptResponse, error) {
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
	Id, err := promptId.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptName := promptui.Prompt{
		Label:   "What is the name of your package",
		Default: "My Package",
	}
	Name, err := promptName.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptDescription := promptui.Prompt{
		Label:   "Provide a description of your package",
		Default: "A package to be used with the platform",
	}
	Description, err := promptDescription.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptStack := promptui.Prompt{
		Label:   "Which stack does your package belong to",
		Default: "instant",
	}
	Stack, err := promptStack.Run()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	promptType := promptui.Select{
		Label: "What type best suites your package",
		Items: []string{"infrastructure", "use-case"},
	}
	index, Type, err := promptType.Run()
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

	var IncludeDevFile bool
	if Dev == "Yes" {
		IncludeDevFile = true
	}

	promptResponse := &GeneratePackagePromptResponse{
		Id:             Id,
		Name:           Name,
		Description:    Description,
		Stack:          Stack,
		Type:           Type,
		IncludeDevFile: IncludeDevFile,
	}

	return promptResponse, nil
}
