package parse

import (
	"path"
	"strings"

	"cli/core"
	"cli/util/slice"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

var (
	ErrNoConfigImage            = errors.New("config file missing field 'image'")
	ErrNoPackages               = errors.New("no packages or custom packages specified")
	ErrUndefinedPackage         = errors.New("packages in command-line not in any of packages, custom-packages, or command-line custom-packages")
	ErrUndefinedProfilePackages = errors.New("packages in profile not in any of packages or custom-packages")
	ErrNoSuchProfile            = errors.New("no such profile")
	ErrNoPackagesInProfile      = errors.New("no packages in profile")
)

func validate(cmd *cobra.Command, config *core.Config) error {
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return errors.Wrap(err, "")
	}

	if config.Image == "" {
		return errors.Wrap(ErrNoConfigImage, "")
	} else if len(config.Packages) == 0 && len(config.CustomPackages) == 0 && len(customPackagePaths) == 0 {
		return errors.Wrap(ErrNoPackages, "")
	}

	packages, err := cmd.Flags().GetStringSlice("name")
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(packages) > 0 {
		err = validateCommandLinePackages(cmd, config)
		if err != nil {
			return err
		}
	}

	profileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		return errors.Wrap(err, "")
	}

	if profileName != "" {
		return validateProfile(cmd, profileName, config)
	}

	return nil
}

func validateCommandLinePackages(cmd *cobra.Command, config *core.Config) error {
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return errors.Wrap(err, "")
	}

	packages, err := cmd.Flags().GetStringSlice("name")
	if err != nil {
		return errors.Wrap(err, "")
	}
	packagesMap := make(map[string]bool)
	for _, pack := range packages {
		packagesMap[pack] = true
		if slice.SliceContains(config.Packages, pack) {
			delete(packagesMap, pack)
			continue
		}

		var hasMatch bool
		for _, customPack := range config.CustomPackages {
			if pack == customPack.Id {
				delete(packagesMap, pack)
				hasMatch = true
				break
			}
		}
		if hasMatch {
			continue
		}

		for _, customPath := range customPackagePaths {
			customPackName := strings.TrimSuffix(path.Base(path.Clean(customPath)), path.Ext(customPath))
			if pack == customPackName {
				delete(packagesMap, pack)
				break
			}
		}
	}

	if len(packagesMap) > 0 {
		for k := range packagesMap {
			return errors.Wrap(ErrUndefinedPackage, k)
		}
	}

	return nil
}

func validateProfile(cmd *cobra.Command, profileName string, config *core.Config) error {
	profilePackagesMap := make(map[string]bool)
	var profile core.Profile
	for _, prof := range config.Profiles {
		if prof.Name == profileName {
			profile = prof
			for _, p := range prof.Packages {
				profilePackagesMap[p] = true
			}

			break
		}
	}

	if profile.Name == "" {
		return errors.Wrap(ErrNoSuchProfile, profileName)
	}
	if len(profilePackagesMap) < 1 {
		return errors.Wrap(ErrNoPackagesInProfile, profileName)
	}

	for _, pack := range profile.Packages {
		if slice.SliceContains(config.Packages, pack) {
			delete(profilePackagesMap, pack)
		}
	}
	for _, pack := range config.CustomPackages {
		if slice.SliceContains(profile.Packages, pack.Id) {
			delete(profilePackagesMap, pack.Id)
		}
	}

	if len(profilePackagesMap) > 0 {
		for k := range profilePackagesMap {
			return errors.Wrap(ErrUndefinedProfilePackages, k)
		}
	}

	return nil
}
