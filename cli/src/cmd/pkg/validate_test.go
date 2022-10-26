package pkg

import (
	"os"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	viperUtil "cli/cmd/util"
	"cli/core"
)

func Test_validate(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		expectedErrorString string
		hookFunc            func(cmd *cobra.Command, config *core.Config)
	}

	configFilePath := wd + "/../../features/unit-test-configs/config-case-5.yml"

	testCases := []cases{
		// case: return ErrNoConfigImage
		{
			expectedErrorString: ErrNoConfigImage.Error(),
			hookFunc:            func(cmd *cobra.Command, config *core.Config) { config.Image = "" },
		},
		// case: return ErrNoPackages
		{
			expectedErrorString: ErrNoPackages.Error(),
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				config.Packages = []string{}
				config.CustomPackages = []core.CustomPackage{}
			},
		},
		// case: no packages specified in config file, but in command-line custom-package, return nil
		{
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				config.Packages = []string{}
				config.CustomPackages = []core.CustomPackage{}

				cmd.Flags().Set("custom-path", "../cares-on-platform")
			},
		},
		// case: return ErrUndefinedProfilePackages
		{
			expectedErrorString: ErrUndefinedProfilePackages.Error(),
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().Set("profile", "dev")
			},
		},
		// case: profile packages not in custom-packages or packages, but in command-line custom-packages, return nil error
		{
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().Set("profile", "dev")
				cmd.Flags().Set("custom-path", "git@github.com:jembi/cares-on-platform.git")
			},
		},
		// case: package specified in command-line that does not exist in config file, expect error
		{
			expectedErrorString: "core: no such command-line package",
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().Set("name", "core")
			},
		},
		// case: no packages specified in command-line, with valid config file, should return nil
		{
			hookFunc: func(cmd *cobra.Command, config *core.Config) {},
		},
	}

	for _, tc := range testCases {
		cmd, config := initCommand(t, configFilePath, tc.hookFunc)

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

	setPackageActionFlags(&cmd)

	hook(&cmd, config)

	return &cmd, config
}
