package parse

import (
	"cli/core"
	"cli/core/state"

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
	if len(profile.EnvFiles) > 0 {
		envViper, err := state.GetEnvironmentVariableViper(profile.EnvFiles)
		if err != nil {
			return nil, err
		}
		envVariables := state.GetEnvVariableString(envViper)
		packageSpec.EnvironmentVariables = append(envVariables, packageSpec.EnvironmentVariables...)
	}

	return &packageSpec, nil
}
