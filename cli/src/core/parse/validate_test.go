package parse

import (
	"os"
	"testing"

	"cli/cmd/flags"
	"cli/core"
	"cli/core/state"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_validate(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		expectedErrorString string
		hookFunc            func(cmd *cobra.Command, config *core.Config)
	}

	configFilePath := wd + "/../../features/unit-test-configs/config-case-4.yml"

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
				cmd.Flags().Set("custom-path", "git@github.com:jembi/disi-on-platform.git")
			},
		},
		// case: no packages specified in command-line, with valid config file, should return nil
		{
			hookFunc: func(cmd *cobra.Command, config *core.Config) {},
		},
		// case: command-line package specified that isn't in config-file, should return nil
		{
			hookFunc: func(cmd *cobra.Command, config *core.Config) {
				cmd.Flags().Set("name", "asdfasdfasdf")
			},
		},
	}

	for _, tc := range testCases {
		cmd, config := initCommand(t, configFilePath, tc.hookFunc)

		err = validate(cmd, config)
		if err != nil {
			require.Equal(t, tc.expectedErrorString, err.Error())
		} else {
			jtest.RequireNil(t, err)
		}
	}
}

func initCommand(t *testing.T, configFilePath string, hook func(cmd *cobra.Command, config *core.Config)) (*cobra.Command, *core.Config) {
	configViper, err := state.SetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(configViper)
	jtest.RequireNil(t, err)

	cmd := &cobra.Command{}
	flags.SetPackageActionFlags(cmd)

	hook(cmd, config)

	return cmd, config
}
