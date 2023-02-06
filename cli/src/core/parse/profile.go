package parse

import (
	"cli/core"
	"cli/core/state"
	"cli/util/slice"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

func getPackageSpecFromProfile(cmd *cobra.Command, config core.Config, packageSpec core.PackageSpec) (*core.PackageSpec, error) {
	profile := core.Profile{}
	profileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	for _, p := range config.Profiles {
		if p.Name == profileName {
			profile = p
			break
		}
	}

	if !cmd.Flags().Changed("dev") && profile.Dev {
		packageSpec.IsDev = profile.Dev
	}

	if !cmd.Flags().Changed("only") && profile.Only {
		packageSpec.IsOnly = profile.Only
	}

	if len(profile.Packages) > 0 {
		packageSpec.Packages = append(profile.Packages, packageSpec.Packages...)
	}

	envVarsMap := make(map[string]string)
	if len(profile.EnvVars) > 0 {
		envVarsMap = slice.AppendUniqueToMapFromSlice(envVarsMap, profile.EnvVars)
	}

	if len(profile.EnvFiles) > 0 {
		envViper, err := state.GetEnvironmentVariableViper(profile.EnvFiles)
		if err != nil {
			return nil, err
		}

		envVariables := state.GetEnvVariableString(envViper)
		envVarsMap = slice.AppendUniqueToMapFromSlice(envVarsMap, envVariables)
	}

	for k, v := range envVarsMap {
		packageSpec.EnvironmentVariables = append(packageSpec.EnvironmentVariables, k+v)
	}

	return &packageSpec, nil
}
