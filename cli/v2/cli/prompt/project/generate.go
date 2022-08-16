package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
)

type GenerateProjectPromptResponse struct {
	ProjectName   string
	ProjectImage  string
	PlatformImage string
}

func GenerateProjectPrompt() GenerateProjectPromptResponse {
	/*
		What is the name of your project:
		What docker image would you like to use with this project: (organisation/project-name)
		What version of platform would you like to use as a base for this project: (latest)
	*/

	path, err := os.Getwd()
	util.LogError(err)

	promptProjectName := promptui.Prompt{
		Label:   "What is the name of your project",
		Default: filepath.Base(path),
	}
	projectName, err := promptProjectName.Run()
	util.LogError(err)

	promptProjectImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this project",
		Default: strcase.ToKebab(fmt.Sprintf("organisation/%v", projectName)),
	}
	projectImage, err := promptProjectImage.Run()
	util.LogError(err)

	promptPlatformImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this project",
		Default: "latest",
	}
	platformImage, err := promptPlatformImage.Run()
	util.LogError(err)

	return GenerateProjectPromptResponse{
		ProjectName:   projectName,
		ProjectImage:  projectImage,
		PlatformImage: platformImage,
	}

}
