package prompt

import (
	"cli/core"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/luno/jettison/errors"
	"github.com/manifoldco/promptui"
)

func GenerateProjectPrompt() (core.Config, error) {
	/*
		What is the name of your project:
		What docker image would you like to use with this project: (organisation/project-name)
		What Platform image is this project based on:
		What is the log path to be used with this project:
		Will this project include custom packages: yes/no
		Would you like to use profiles in this project:
	*/

	path, err := os.Getwd()
	if err != nil {
		return core.Config{}, errors.Wrap(err, "")
	}

	promptProjectName := promptui.Prompt{
		Label:   "What is the name of your project",
		Default: filepath.Base(path),
	}
	projectName, err := promptProjectName.Run()
	if err != nil {
		return core.Config{}, errors.Wrap(err, "")
	}

	promptProjectImage := promptui.Prompt{
		Label:   "What docker image would you like to use with this project",
		Default: strcase.ToKebab(fmt.Sprintf("organisation/%v", projectName)),
	}
	projectImage, err := promptProjectImage.Run()
	if err != nil {
		return core.Config{}, errors.Wrap(err, "")
	}

	promptPlatformImage := promptui.Prompt{
		Label:   "What Platform image is this project based on",
		Default: "jembi/platform",
	}
	platformImage, err := promptPlatformImage.Run()
	if err != nil {
		return core.Config{}, errors.Wrap(err, "")
	}

	promptLogPath := promptui.Prompt{
		Label:   "What is the log path to be used with this project",
		Default: "/tmp/logs",
	}
	logPath, err := promptLogPath.Run()
	if err != nil {
		return core.Config{}, errors.Wrap(err, "")
	}

	promptPackages := promptui.Select{
		Label: "Will this project include custom packages",
		Items: []string{"Yes", "No"},
	}
	index, withCustomPackages, err := promptPackages.Run()
	if err != nil || index == -1 {
		return core.Config{}, errors.Wrap(err, "")
	}

	var customPackages []core.CustomPackage
	if withCustomPackages == "Yes" {
		customPackages = append(customPackages, core.CustomPackage{
			Id:   "<<custom-package-id>>",
			Path: "<<custom-package-path>>",
		})
	}

	promptProfiles := promptui.Select{
		Label: "Will this project include project profiles",
		Items: []string{"Yes", "No"},
	}
	index, withProfiles, err := promptProfiles.Run()
	if err != nil || index == -1 {
		return core.Config{}, errors.Wrap(err, "")
	}

	var profiles []core.Profile
	if withProfiles == "Yes" {
		profiles = append(profiles, core.Profile{
			Name:     "<<profile-name>>",
			EnvFiles: []string{"<<env-file-1>>", "<<env-file-2>>"},
			EnvVars: []string{"<<env-var-1>>", "<<env-var-2>>"},
			Packages: []string{"<<profile-package-id-1>>", "<<profile-package-id-2>>"},
		})
	}

	return core.Config{
		Image:          projectImage,
		ProjectName:    projectName,
		PlatformImage:  platformImage,
		LogPath:        logPath,
		Packages:       []string{"<<package-id>>"},
		CustomPackages: customPackages,
		Profiles:       profiles,
	}, nil
}
