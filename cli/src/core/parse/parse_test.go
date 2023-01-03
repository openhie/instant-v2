package parse

import (
	"cli/util/slice"
	"os"
	"strings"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestParseAndPrepareLaunch(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		configFilePath       string
		expectedErrorString  string
		duplicatedEnvVarName string
		expectedEnvVarValue  string
		hookFunc             func(cmd *cobra.Command)
	}

	testCases := []cases{
		// case: command line env file env vars must overwrite profile env file env vars
		{
			configFilePath:       wd + "/../../features/unit-test-configs/config-case-1.yml",
			expectedErrorString:  ".env.none: no such file or directory",
			duplicatedEnvVarName: "FIRST_ENV_VAR",
			expectedEnvVarValue:  "not_number_one",
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("config", wd+"/../../features/unit-test-configs/config-case-1.yml")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.three")
				jtest.RequireNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		cmd, _ := loadCmdAndConfig(t, tc.configFilePath, tc.hookFunc)

		pSpec, _, err := ParseAndPrepareLaunch(cmd)
		jtest.RequireNil(t, err)

		require.Equal(t, 1, substringInstancesInSlice(pSpec.EnvironmentVariables, tc.duplicatedEnvVarName))
		require.Equal(t, true, slice.SliceContains(pSpec.EnvironmentVariables, tc.duplicatedEnvVarName+"="+tc.expectedEnvVarValue))
	}
}

func substringInstancesInSlice(slice []string, element string) int {
	var instances int
	for _, s := range slice {
		if strings.Contains(s, element) {
			instances++
		}
	}

	return instances
}
