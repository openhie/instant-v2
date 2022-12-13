package parse

import (
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"cli/cmd/flags"
	"cli/core"
	"cli/core/state"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func copyFile(src, dst string) error {
	// Open the source file
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	// Create the destination file
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	return nil

}

func TestGetPackageSpecFromParams(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		configFilePath string
		hookFunc       func(cmd *cobra.Command)
		wantSpecMatch  bool
		packageSpec    *core.PackageSpec
		errorString    string
	}

	testCases := []cases{
		// case: match packageSpec
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-3.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().StringSlice("env-file", []string{""}, "")

				cmd.Flags().Set("name", "pack-1")
				cmd.Flags().Set("name", "pack-2")

				cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.one")
				cmd.Flags().Set("env-file", wd+"/../../features/test-conf/.env.two")

				cmd.Flags().Set("custom-path", "disi-on-platform")

				cmd.Flags().Set("dev", "true")
				cmd.Flags().Set("only", "true")
			},
			wantSpecMatch: true,
			packageSpec: &core.PackageSpec{
				Packages: []string{"pack-1", "pack-2"},
				CustomPackages: []core.CustomPackage{
					{
						Id:   "disi-on-platform",
						Path: "git@github.com:jembi/disi-on-platform.git",
					},
				},
				EnvironmentVariables: []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"},
				IsDev:                true,
				IsOnly:               true,
			},
		},
		// case: return error from not finding env file
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-3.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().StringSlice("env-file", []string{""}, "")

				cmd.Flags().Set("name", "pack-1")
				cmd.Flags().Set("env-file", wd+"/../../features/test-conf/awlikdeuh")
			},
			errorString: "no such file or directory",
		},
		// case: return no error when not specifying an env-file
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-3.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().StringSlice("env-file", []string{""}, "")

				cmd.Flags().Set("name", "pack-1")
			},
		},
		// case: place .env file in main dir, but don't use its env vars
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-3.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().StringSlice("env-file", []string{""}, "")

				cmd.Flags().Set("name", "pack-1")

				cmd.Flags().Set("only", "true")
			},
			wantSpecMatch: true,
			packageSpec: &core.PackageSpec{
				Packages: []string{"pack-1"},
				IsOnly:   true,
			},
		},
	}

	for _, tc := range testCases {
		defer os.Remove(wd + "/../../.env")
		err = copyFile(wd+"/../../features/test-conf/.env.test", wd+"/../../.env")
		jtest.RequireNil(t, err)

		cmd, config := loadCmdAndConfig(t, tc.configFilePath, tc.hookFunc)

		pSpec, err := GetPackageSpecFromParams(cmd, config)
		if tc.errorString != "" && !strings.Contains(err.Error(), tc.errorString) {
			t.FailNow()
		} else if tc.errorString == "" {
			jtest.RequireNil(t, err)
		}

		if tc.wantSpecMatch {
			sort.Slice(pSpec.EnvironmentVariables, func(i, j int) bool {
				return strings.Contains(pSpec.EnvironmentVariables[i], "FIRST_ENV_VAR")
			})

			if !assert.Equal(t, tc.packageSpec, pSpec) {
				t.FailNow()
			}
		}
	}
}

func loadCmdAndConfig(t *testing.T, configFilePath string, hookFunc func(cmd *cobra.Command)) (*cobra.Command, *core.Config) {
	configViper, err := state.SetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(configViper)
	jtest.RequireNil(t, err)

	cmd := &cobra.Command{}
	flags.SetPackageActionFlags(cmd)

	hookFunc(cmd)

	return cmd, config
}

var (
	expectedCustomPackages = []core.CustomPackage{
		{
			Id:   "custom-package-1",
			Path: "path-to-1",
		},
		{
			Id:   "custom-package-2",
			Path: "path-to-2",
		},
	}
)

func Test_getCustomPackages(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	configViper, err := state.SetConfigViper(wd + "/../../features/unit-test-configs/config-case-4.yml")
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(configViper)
	jtest.RequireNil(t, err)

	gotCustomPackages := getCustomPackages(config, []string{"path-to-1", "path-to-2"})

	require.Equal(t, expectedCustomPackages, gotCustomPackages)
}
