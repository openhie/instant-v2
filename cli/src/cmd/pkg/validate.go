package pkg

import (
	"cli/core"
	"cli/util"
	"path"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

var (
	ErrNoConfigImage            = errors.New("config file missing field 'image'")
	ErrNoPackages               = errors.New("no packages or custom packages specified")
	ErrUndefinedProfilePackages = errors.New("packages in profile not in any of packages, custom-packages, or command-line custom-packages")
)

func validate(cmd *cobra.Command, config *core.Config) error {
	customPackagePaths, err := cmd.Flags().GetStringSlice("custom-path")
	if err != nil && !strings.Contains(err.Error(), "flag accessed but not defined") {
		return errors.Wrap(err, "")
	}

	if config.Image == "" {
		return errors.Wrap(ErrNoConfigImage, "")
	} else if len(config.Packages) == 0 && len(config.CustomPackages) == 0 && len(customPackagePaths) == 0 {
		return errors.Wrap(ErrNoPackages, "")
	}

	profile, err := cmd.Flags().GetString("profile")
	if err != nil && !strings.Contains(err.Error(), "flag accessed but not defined") {
		return errors.Wrap(err, "")
	} else {
		err = validateProfile(cmd, config, profile)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateProfile(cmd *cobra.Command, config *core.Config, profile string) error {
	profilePackagesMap := make(map[string]bool)
	for _, pack := range config.Profiles {
		if pack.Name == profile {
			for _, p := range pack.Packages {
				profilePackagesMap[p] = true
			}
		}
	}

	for _, profile := range config.Profiles {
		for _, pack := range profile.Packages {
			if util.SliceContains(config.Packages, pack) {
				delete(profilePackagesMap, pack)
			}
		}
		for _, pack := range config.CustomPackages {
			if !util.SliceContains(profile.Packages, pack.Id) {
				delete(profilePackagesMap, pack.Id)
			}
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
