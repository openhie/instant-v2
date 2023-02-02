package parse

import (
	"os"
	"sort"
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
		duplicatedEnvVarName string
		expectedEnvVars      []string
		hookFunc             func(cmd *cobra.Command)
	}

	testCases := []cases{
		// case: command-line env file env vars must overwrite profile env file env vars
		{
			configFilePath:       wd + "/../../features/unit-test-configs/config-case-1.yml",
			duplicatedEnvVarName: "FIRST_ENV_VAR",
			expectedEnvVars:      []string{"FIRST_ENV_VAR=not_number_one", "SECOND_ENV_VAR=number_two"},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("config", wd+"/../../features/unit-test-configs/config-case-1.yml")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.three")
				jtest.RequireNil(t, err)
			},
		},
		// case: command-line env-vars must overwrite profile env file env vars
		{
			configFilePath:       wd + "/../../features/unit-test-configs/config-case-1.yml",
			duplicatedEnvVarName: "FIRST_ENV_VAR",
			expectedEnvVars:      []string{"FIRST_ENV_VAR=command_line_value", "SECOND_ENV_VAR=number_two"},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("config", wd+"/../../features/unit-test-configs/config-case-1.yml")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-var", "FIRST_ENV_VAR=command_line_value")
				jtest.RequireNil(t, err)
			},
		},
		// case: command-line env-vars must overwrite command-line env file env vars
		{
			configFilePath:       wd + "/../../features/unit-test-configs/config-case-1.yml",
			duplicatedEnvVarName: "FIRST_ENV_VAR",
			expectedEnvVars:      []string{"FIRST_ENV_VAR=command_line_value"},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("config", wd+"/../../features/unit-test-configs/config-case-1.yml")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.three")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-var", "FIRST_ENV_VAR=command_line_value")
				jtest.RequireNil(t, err)
			},
		},
		// case: command-line env-vars must overwrite profile env file env vars and command-line env file env vars
		{
			configFilePath:       wd + "/../../features/unit-test-configs/config-case-1.yml",
			duplicatedEnvVarName: "FIRST_ENV_VAR",
			expectedEnvVars:      []string{"FIRST_ENV_VAR=command_line_value", "SECOND_ENV_VAR=number_two"},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("config", wd+"/../../features/unit-test-configs/config-case-1.yml")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.three")
				jtest.RequireNil(t, err)

				err = cmd.Flags().Set("env-var", "FIRST_ENV_VAR=command_line_value")
				jtest.RequireNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		cmd, _ := loadCmdAndConfig(t, tc.configFilePath, tc.hookFunc)

		pSpec, _, err := ParseAndPrepareLaunch(cmd)
		jtest.RequireNil(t, err)

		require.Equal(t, 1, substringInstancesInSlice(pSpec.EnvironmentVariables, tc.duplicatedEnvVarName))

		sort.Slice(pSpec.EnvironmentVariables, func(i, j int) bool {
			return strings.Contains(pSpec.EnvironmentVariables[i], "FIRST_ENV_VAR")
		})

		require.Equal(t, tc.expectedEnvVars, pSpec.EnvironmentVariables)
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
