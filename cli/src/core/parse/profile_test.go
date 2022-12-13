package parse

import (
	"cli/core"
	"context"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGetPackageSpecFromProfile(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		configFilePath      string
		expectedErrorString string
		expectedConfig      *core.PackageSpec
		hookFunc            func(cmd *cobra.Command)
	}

	// TODO: throw error if specifying non-existant profile
	testCases := []cases{
		// case: return error from non-existant env file directory
		{
			configFilePath:      wd + "/../../features/unit-test-configs/config-case-5.yml",
			expectedErrorString: ".env.none: no such file or directory",
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "bad-env-file-path")
				jtest.RequireNil(t, err)
			},
		},
		// case: assert dev profile from config-case-5.yml parsed properly
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-5.yml",
			expectedConfig: &core.PackageSpec{
				EnvironmentVariables: []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"},
				Packages:             []string{"dashboard-visualiser-jsreport", "disi-on-platform"},
				IsDev:                true,
			},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev")
				jtest.RequireNil(t, err)
			},
		},
		// case: assert only profile from config-case-5.yml parsed properly
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-5.yml",
			expectedConfig: &core.PackageSpec{
				EnvironmentVariables: []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"},
				Packages:             []string{"dashboard-visualiser-jsreport", "disi-on-platform", "core"},
				IsOnly:               true,
			},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "only")
				jtest.RequireNil(t, err)
			},
		},
		// case: assert dev-and-only profile from config-case-5.yml parsed properly
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-5.yml",
			expectedConfig: &core.PackageSpec{
				EnvironmentVariables: []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"},
				Packages:             []string{"core"},
				IsDev:                true,
				IsOnly:               true,
			},
			hookFunc: func(cmd *cobra.Command) {
				err = cmd.Flags().Set("profile", "dev-and-only")
				jtest.RequireNil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		cmd, config := loadCmdAndConfig(t, tc.configFilePath, tc.hookFunc)

		pSpec, err := GetPackageSpecFromProfile(cmd, *config, core.PackageSpec{})
		if err != nil {
			require.Equal(t, strings.Contains(err.Error(), tc.expectedErrorString), true)
		} else if err != nil && !strings.Contains(err.Error(), tc.expectedErrorString) {
			log.Error(context.Background(), err)
			t.FailNow()
		} else if tc.expectedConfig != nil {
			sort.Slice(pSpec.EnvironmentVariables, func(i, j int) bool {
				return strings.Contains(pSpec.EnvironmentVariables[i], "FIRST_ENV_VAR")
			})

			require.Equal(t, tc.expectedConfig, pSpec)
		}
	}
}
