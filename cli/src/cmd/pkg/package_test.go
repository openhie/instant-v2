package pkg

import (
	"os"
	"testing"

	viperUtil "cli/cmd/util"
	"cli/core"

	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var configCaseOne = core.Config{
	Image:    "jembi/go-cli-test-image",
	LogPath:  "/tmp/logs",
	Packages: []string{"client", "dashboard-visualiser-jsreport"},
	CustomPackages: []core.CustomPackage{
		{
			Id:          "disi-on-platform",
			Path:        "git@github.com:jembi/disi-on-platform.git",
			SshKey:      "./id_rsa",
			SshPassword: "./id_rsa_password.txt",
		},
	},
	Profiles: []core.Profile{
		{
			Name:     "dev",
			EnvFiles: []string{"../test-conf/.env.test"},
			Packages: []string{"dashboard-visualiser-jsreport", "disi-on-platform"},
			Dev:      true,
		},
	},
}

func Test_unmarshalConfig(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	configViper, err := viperUtil.GetConfigViper(wd + "/../../features/unit-test-configs/config-case-1.yml")
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	if !assert.Equal(t, configCaseOne, *config) {
		t.FailNow()
	}
}

func Test_loadInProfileParams(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		profileName         string
		boolFlagName        string
		configFilePath      string
		expectedErrorString string
	}

	testCases := []cases{
		// case: return error from conflicting command-line --dev flag and 'dev: false' config.yml profile flag
		{
			profileName:         "non-dev",
			boolFlagName:        "dev",
			configFilePath:      wd + "/../../features/unit-test-configs/config-case-2.yml",
			expectedErrorString: ErrConflictingDevFlag.Error(),
		},
		// case: return error from conflicting command-line --only flag and 'only: false' config.yml profile flag
		{
			profileName:         "non-only",
			boolFlagName:        "only",
			configFilePath:      wd + "/../../features/unit-test-configs/config-case-2.yml",
			expectedErrorString: ErrConflictingOnlyFlag.Error(),
		},
		// case: return error from non-existant env file directory
		{
			profileName:         "bad-env-file-path",
			boolFlagName:        "",
			configFilePath:      wd + "/../../features/unit-test-configs/config-case-3.yml",
			expectedErrorString: "stat ./features/test-conf/.env.tests: no such file or directory",
		},
	}

	for _, tc := range testCases {
		cmd, config := setupLoadInProfileParams(t, tc.configFilePath)

		setupBoolFlags(t, cmd, tc.boolFlagName)
		cmd.Flags().String("profile", tc.profileName, "")

		_, err = loadInProfileParams(cmd, *config, core.PackageSpec{})
		if !assert.Equal(t, tc.expectedErrorString, err.Error()) {
			t.FailNow()
		}
	}

	// case: load in environment variables from more than one env file
	cmd, config := setupLoadInProfileParams(t, wd+"/../../features/unit-test-configs/config-case-2.yml")

	cmd.Flags().String("profile", "non-only", "")

	packageSpec, err := loadInProfileParams(cmd, *config, core.PackageSpec{})
	jtest.RequireNil(t, err)

	if !assert.Equal(t, packageSpec.EnvironmentVariables, []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"}) {
		t.FailNow()
	}
}

func setupLoadInProfileParams(t *testing.T, configFilePath string) (*cobra.Command, *core.Config) {
	configViper, err := viperUtil.GetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	return &cobra.Command{}, config
}

func setupBoolFlags(t *testing.T, cmd *cobra.Command, boolFlagName string) {
	cmd.Flags().Bool(boolFlagName, false, "")
	err := cmd.Flags().Set(boolFlagName, "true")
	jtest.RequireNil(t, err)
}
