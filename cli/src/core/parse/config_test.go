package parse

import (
	"os"
	"testing"

	"cli/core"
	coreConfig "cli/core/state"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
)

func Test_appendTag(t *testing.T) {
	type cases struct {
		config        *core.Config
		wantImageName string
	}

	testCases := []cases{
		{
			config: &core.Config{
				Image: "docker/image",
			},
			wantImageName: "docker/image:latest",
		},
		{
			config: &core.Config{
				Image: "docker/image:latest",
			},
			wantImageName: "docker/image:latest",
		},
		{
			config: &core.Config{
				Image: "docker/image:1.0.0",
			},
			wantImageName: "docker/image:1.0.0",
		},
	}

	for _, tc := range testCases {
		appendTag(tc.config)
		require.Equal(t, tc.wantImageName, tc.config.Image)
	}
}

func Test_unmarshalConfig(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		configPath     string
		expectedConfig core.Config
		errString      string
	}

	testCases := []cases{
		// case: match configCaseOne
		{
			configPath: wd + "/../../features/unit-test-configs/config-case-1.yml",
			expectedConfig: core.Config{
				ProjectName: "test-project",
				Image:         "jembi/go-cli-test-image",
				LogPath:       "/tmp/logs",
				PlatformImage: "jembi/platform:latest",
				Packages:      []string{"client", "dashboard-visualiser-jsreport"},
				CustomPackages: []core.CustomPackage{
					{
						Id:   "disi-on-platform",
						Path: "git@github.com:jembi/disi-on-platform.git",
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
			},
		},
		// case: return invalid config file syntax error
		{
			configPath: wd + "/../../features/unit-test-configs/config-case-2.yml",
			errString:  ErrInvalidConfigFileSyntax.Error(),
		},
	}

	for _, testCase := range testCases {
		configViper, err := coreConfig.SetConfigViper(testCase.configPath)
		jtest.RequireNil(t, err)

		config, err := unmarshalConfig(configViper)
		if testCase.errString == "" {
			jtest.RequireNil(t, err)
			require.Equal(t, testCase.expectedConfig, *config)
		} else {
			require.Equal(t, ErrInvalidConfigFileSyntax.Error(), err.Error())
		}
	}
}
