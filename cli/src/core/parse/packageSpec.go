package parse

import (
	"cli/core"
	"cli/core/state"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

// Match custom packages passed through from the command line to custom
// packages specified in the config file (if they exist in the config file),
// otherwise append the custom package path but not ID.
func parseCustomPackageFromPath(config *core.Config, customPackagePaths []string) []core.CustomPackage {
	var customPackages []core.CustomPackage
	for _, customPackagePath := range customPackagePaths {
		var customPackage core.CustomPackage

		for _, configCustomPackage := range config.CustomPackages {
			if customPackagePath == configCustomPackage.Id || customPackagePath == configCustomPackage.Path {
				customPackage = configCustomPackage
				break
			}
		}
		if customPackage.Id == "" {
			customPackage = core.CustomPackage{
				Path: customPackagePath,
			}
		}

		customPackages = append(customPackages, customPackage)
	}

	return customPackages
}

func GetPackageSpecFromParams(cmd *cobra.Command, config *core.Config) (*core.PackageSpec, error) {
	packageSpec := core.PackageSpec{}

	packageNames, err := cmd.Flags().GetStringSlice("name")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	isDev, err := cmd.Flags().GetBool("dev")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	isOnly, err := cmd.Flags().GetBool("only")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var envVariables []string
	if cmd.Flags().Changed("env-file") {
		envFiles, err := cmd.Flags().GetStringSlice("env-file")
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		envViper, err := state.GetEnvironmentVariableViper(envFiles)
		if err != nil {
			return nil, err
		}
		envVariables = state.GetEnvVariableString(envViper)
	}

	customPackages := parseCustomPackageFromPath(config, customPackagePaths)

	packageSpec = core.PackageSpec{
		Packages:             packageNames,
		CustomPackages:       customPackages,
		EnvironmentVariables: envVariables,
		IsDev:                isDev,
		IsOnly:               isOnly,
		DeployCommand:        cmd.Use,
	}

	return &packageSpec, nil
}
