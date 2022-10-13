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

	// case: return error from conflicting command-line --dev flag and dev: false config.yml profile flag
	cmd, config := setupLoadInProfileParams(t, wd+"/../../features/unit-test-configs/config-case-2.yml")

	cmd.Flags().Bool("dev", false, "")
	err = cmd.Flags().Set("dev", "true")
	jtest.RequireNil(t, err)

	cmd.Flags().String("profile", "non-dev", "")

	_, err = loadInProfileParams(cmd, *config, core.PackageSpec{})
	if !assert.Equal(t, ErrConflictingDevFlag.Error(), err.Error()) {
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
			EnvFiles: []string{"./features/test-conf/.env.tests"},
			Packages: []string{"dashboard-visualiser-jsreport", "disi-on-platform"},
			Dev:      true,
		},
	},
}
