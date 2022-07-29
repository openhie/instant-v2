package pkg

import (
	"testing"
)

func Test_executeCommand(t *testing.T) {
	tests := []struct {
		name           string
		CustomOptions  customOption
		deployCommands []string
	}{
		{
			name: "Test case assert startupPackages",
			CustomOptions: customOption{
				startupAction:   "down",
				startupPackages: []string{"core", "elastic-analytics"},
				imageVersion:    "latest",
				targetLauncher:  "docker",
			},
			deployCommands: []string{"down", "core", "elastic-analytics", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert envVarFileLocation",
			CustomOptions: customOption{
				startupAction:      "init",
				envVarFileLocation: "./usr/bin",
				imageVersion:       "latest",
				targetLauncher:     "k8s",
			},
			deployCommands: []string{"init", "--env-file=./usr/bin", "--image-version=latest", "-t=k8s"},
		},
		{
			name: "Test case assert envVars",
			CustomOptions: customOption{
				startupAction:  "up",
				envVars:        []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				imageVersion:   "latest",
				targetLauncher: "docker",
			},
			deployCommands: []string{"up", "-e=NODE_ENV=DEV", "-e=DOMAIN_NAME=instant.com", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert customPackageFileLocations",
			CustomOptions: customOption{
				startupAction:              "init",
				customPackageFileLocations: []string{"./local/cPack"},
				imageVersion:               "latest",
				targetLauncher:             "docker",
			},
			deployCommands: []string{"init", "-c=./local/cPack", "--image-version=latest", "-t=docker"},
		},
		{
			name: "Test case assert dev and only flags",
			CustomOptions: customOption{
				startupAction:  "destroy",
				onlyFlag:       true,
				imageVersion:   "v1.02b",
				targetLauncher: "k8s",
				devMode:        true,
			},
			deployCommands: []string{"destroy", "--only", "--dev", "--image-version=v1.02b", "-t=k8s"},
		},
		{
			name: "Test case assert all fields",
			CustomOptions: customOption{
				startupAction:              "init",
				startupPackages:            []string{"hmis", "mcsd"},
				envVarFileLocation:         "./home/bin",
				envVars:                    []string{"NODE_ENV=DEV", "DOMAIN_NAME=instant.com"},
				customPackageFileLocations: []string{"./usr/local/cPack"},
				onlyFlag:                   true,
				imageVersion:               "v1.03a",
				targetLauncher:             "k8s",
				devMode:                    true,
			},
			deployCommands: []string{"init", "hmis", "mcsd", "--env-file=./home/bin", "-e=NODE_ENV=DEV",
				"-e=DOMAIN_NAME=instant.com", "-c=./usr/local/cPack", "--only", "--dev", "--image-version=v1.03a", "-t=k8s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customOptions = tt.CustomOptions
			runDeployCommand = func(startupCommands []string) error {
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
