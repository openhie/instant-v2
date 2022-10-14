package pkg

import (
	viperUtil "cli/cmd/util"
	"cli/core"
	"os"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_validate(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		preperationFunc     func(t *testing.T, configFilePath string, hook func(cmd *cobra.Command, config *core.Config)) (*cobra.Command, *core.Config)
		expectedErrorString string
		hookFunc            func(cmd *cobra.Command, config *core.Config)
	}

	configFilePath := wd + "/../../features/unit-test-configs/config-case-5.yml"

	testCases := []cases{
		// case: return ErrNoConfigImage
		{
			preperationFunc:     initCommand,
			expectedErrorString: ErrNoConfigImage.Error(),
			hookFunc:            func(cmd *cobra.Command, config *core.Config) { config.Image = "" },
		},
		// case: return ErrNoPackages
		{
			preperationFunc:     initCommand,
			expectedErrorString: ErrNoPackages.Error(),
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				config.Packages = []string{}
				config.CustomPackages = []core.CustomPackage{}
			},
		},
		// case: no packages specified in config file, but in command-line custom-package, return nil
		{
			preperationFunc: initCommand,
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				config.Packages = []string{}
				config.CustomPackages = []core.CustomPackage{}

				cmd.Flags().StringSlice("custom-path", []string{"../cares-on-platform"}, "")
			},
		},
		// case: return ErrUndefinedProfilePackages
		{
			preperationFunc:     initCommand,
			expectedErrorString: ErrUndefinedProfilePackages.Error(),
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().StringSlice("profile", []string{"dev"}, "")
			},
		},
		// case: profile packages not in custom-packages or packages, but in command-line custom-packages, return nil error
		{
			preperationFunc: initCommand,
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().StringSlice("profile", []string{"dev"}, "")
				cmd.Flags().StringSlice("custom-path", []string{"git@github.com:jembi/cares-on-platform.git"}, "")
			},
		},
	}

	for _, tc := range testCases {
		cmd, config := tc.preperationFunc(t, configFilePath, tc.hookFunc)

		err = validate(cmd, config)
		if err != nil {
			if !assert.Equal(t, tc.expectedErrorString, err.Error()) {
				t.FailNow()
			}
		} else {
			if !assert.Equal(t, tc.expectedErrorString, "") {
				t.FailNow()
			}
		}
	}
}

func initCommand(t *testing.T, configFilePath string, hook func(cmd *cobra.Command, config *core.Config)) (*cobra.Command, *core.Config) {
	configViper, err := viperUtil.GetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	cmd := cobra.Command{}

	hook(&cmd, config)

	return &cmd, config
}
