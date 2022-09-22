package prompts

import (
	"testing"

	"ohiecli/config"
	"ohiecli/docker"
)

func Test_executeCommand(t *testing.T) {
	tests := []struct {
		name           string
		CustomOptions  config.CustomOption
		deployCommands []string
	}{
		{
			name: "Test case assert startupPackages",
			CustomOptions: config.CustomOption{
				StartupAction:   "down",
				StartupPackages: []string{"core", "elastic-analytics"},
				ImageVersion:    "latest",
				TargetLauncher:  "docker",
			},
			deployCommands: []string{"down", "core", "elastic-analytics", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert envVarFileLocation",
			CustomOptions: config.CustomOption{
				StartupAction:      "init",
				EnvVarFileLocation: "./usr/bin",
				ImageVersion:       "latest",
				TargetLauncher:     "k8s",
			},
			deployCommands: []string{"init", "--env-file=./usr/bin", "--image-version=latest", "-t=k8s"},
		},
		{
			name: "Test case assert envVars",
			CustomOptions: config.CustomOption{
				StartupAction:  "up",
				EnvVars:        []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				ImageVersion:   "latest",
				TargetLauncher: "docker",
			},
			deployCommands: []string{"up", "-e=NODE_ENV=DEV", "-e=DOMAIN_NAME=instant.com", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert customPackageFileLocations",
			CustomOptions: config.CustomOption{
				StartupAction:              "init",
				CustomPackageFileLocations: []string{"./local/cPack"},
				ImageVersion:               "latest",
				TargetLauncher:             "docker",
			},
			deployCommands: []string{"init", "-c=./local/cPack", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert dev and only flags",
			CustomOptions: config.CustomOption{
				StartupAction:  "destroy",
				OnlyFlag:       true,
				ImageVersion:   "v1.02b",
				TargetLauncher: "k8s",
				DevMode:        true,
			},
			deployCommands: []string{"destroy", "--only", "--dev", "--image-version=v1.02b", "-t=k8s"},
		},
		{
			name: "Test case assert all fields",
			CustomOptions: config.CustomOption{
				StartupAction:              "init",
				StartupPackages:            []string{"hmis", "mcsd"},
				EnvVarFileLocation:         "./home/bin",
				EnvVars:                    []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				CustomPackageFileLocations: []string{"./usr/local/cPack"},
				OnlyFlag:                   true,
				ImageVersion:               "v1.03a",
				TargetLauncher:             "k8s",
				DevMode:                    true,
			},
			deployCommands: []string{"init", "hmis", "mcsd", "--env-file=./home/bin", "-e=NODE_ENV=DEV",
				"-e=DOMAIN_NAME=instant.com", "-c=./usr/local/cPack", "--only", "--dev", "--image-version=v1.03a", "-t=k8s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CustomOptions = tt.CustomOptions
			docker.RunDeployCommand = func(startupCommands []string) error {
				return nil
			}

			executeCommand()
			for i, dc := range DeployCommands {
				if dc != tt.deployCommands[i] {
					t.Errorf("DeployCommands variable error, got = %v, expected %v", dc, tt.deployCommands[i])
				}
			}
		})
	}
}
