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
	ErrUndefinedProfilePackages = errors.New("packages in profile not in any of packages, custom-packages, or command-line custom-packages")
	ErrNoSuchProfile            = errors.New("no such profile")
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

	profileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		return errors.Wrap(err, "")
	}

	if profileName != "" {
		return validateProfile(cmd, profileName, config)
	}

	return nil
}

func validateProfile(cmd *cobra.Command, profileName string, config *core.Config) error {
	profilePackagesMap := make(map[string]bool)
	var profile core.Profile
	for _, prof := range config.Profiles {
		if prof.Name == profileName {
			for _, p := range prof.Packages {
				profilePackagesMap[p] = true
				profile = prof
			}
			
			break
		}
	}

	if len(profilePackagesMap) < 1 {
		return errors.Wrap(ErrNoSuchProfile, profileName)
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

	customPackages, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, cp := range customPackages {
		packName := strings.TrimSuffix(path.Base(path.Clean(cp)), path.Ext(cp))
		if profilePackagesMap[packName] {
			delete(profilePackagesMap, packName)
		}
	}

	if len(profilePackagesMap) > 0 {
		return errors.Wrap(ErrUndefinedProfilePackages, "")
	}

	return nil
}
