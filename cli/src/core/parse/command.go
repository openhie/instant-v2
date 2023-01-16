package parse

import (
	"path"
	"strings"

	"cli/core"
	"cli/util/slice"
)

func GetCustomPackageName(customPackage core.CustomPackage) string {
	if customPackage.Id != "" {
		return customPackage.Id
	}
	return strings.TrimSuffix(path.Base(path.Clean(customPackage.Path)), path.Ext(customPackage.Path))
}

func GetInstantCommand(packageSpec core.PackageSpec) []string {
	instantCommand := []string{packageSpec.DeployCommand, "-t", "swarm"}

	if packageSpec.IsDev {
		instantCommand = append(instantCommand, "--dev")
	}

	if packageSpec.IsOnly {
		instantCommand = append(instantCommand, "--only")
	}

	instantCommand = append(instantCommand, packageSpec.Packages...)

	for _, customPackage := range packageSpec.CustomPackages {
		customPackageName := GetCustomPackageName(customPackage)
		if !slice.SliceContains(packageSpec.Packages, customPackageName) {
			instantCommand = append(instantCommand, customPackageName)
		}
	}

	return instantCommand
}
