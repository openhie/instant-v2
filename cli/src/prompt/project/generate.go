package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/luno/jettison/errors"
	"github.com/manifoldco/promptui"
)

type GenerateProjectPromptResponse struct {
	ProjectName   string
	ProjectImage  string
	PlatformImage string
}

func GenerateProjectPrompt() (GenerateProjectPromptResponse, error) {
	/*
		What is the name of your project:
		What docker image would you like to use with this project: (organisation/project-name)
	*/

	path, err := os.Getwd()
	if err != nil {
		return GenerateProjectPromptResponse{}, errors.Wrap(err, "")
	}

	promptProjectName := promptui.Prompt{
		Label:   "What is the name of your project",
		Default: filepath.Base(path),
	}
	projectName, err := promptProjectName.Run()
	if err != nil {
		return GenerateProjectPromptResponse{}, errors.Wrap(err, "")
	}

	promptProjectImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this project",
		Default: strcase.ToKebab(fmt.Sprintf("organisation/%v", projectName)),
	}
	projectImage, err := promptProjectImage.Run()
	if err != nil {
		return GenerateProjectPromptResponse{}, errors.Wrap(err, "")
	}

	return GenerateProjectPromptResponse{
		ProjectName:  projectName,
		ProjectImage: projectImage,
	}, nil
}
