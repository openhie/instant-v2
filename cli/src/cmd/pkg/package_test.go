package pkg

import (
	"os"
	"testing"

	viperUtil "cli/cmd/util"
	"cli/core"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/assert"
)

func Test_getConfigFromParams(t *testing.T) {
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
