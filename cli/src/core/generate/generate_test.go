package generate

import (
	"bufio"
	"bytes"
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"

	"cli/core"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
)

func Test_createFileFromTemplate(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		src                 string
		dst                 string
		generatePackageSpec core.GeneratePackageSpec
		pathToExpectedFile  string
	}

	generatePackSpec := core.GeneratePackageSpec{
		Id:             "test-package",
		Name:           "Test Package",
		Stack:          "test-stack",
		Image:          "test/image",
		Description:    "A package for testing",
		Type:           "infrastructure",
		IncludeDevFile: true,
		TargetPort:     "1",
		PublishedPort:  "2",
	}

	testCases := []cases{
		// case: assert swarm.sh is created from template as expected
		{
			src:                 "swarm.sh",
			dst:                 filepath.Join(wd, "test-package"),
			generatePackageSpec: generatePackSpec,
			pathToExpectedFile:  filepath.Join(wd, "..", "..", "features", "test-package", "swarm.sh"),
		},
		// case: assert package-metadata.json is created from template as expected
		{
			src:                 "package-metadata.json",
			dst:                 filepath.Join(wd, "test-package"),
			generatePackageSpec: generatePackSpec,
			pathToExpectedFile:  filepath.Join(wd, "..", "..", "features", "test-package", "package-metadata.json"),
		},
		// case: assert docker-compose.yml is created from template as expected
		{
			src:                 "docker-compose.yml",
			dst:                 filepath.Join(wd, "test-package"),
			generatePackageSpec: generatePackSpec,
			pathToExpectedFile:  filepath.Join(wd, "..", "..", "features", "test-package", "docker-compose.yml"),
		},
		// case: assert docker-compose.dev.yml is created from template as expected
		{
			src:                 "docker-compose.dev.yml",
			dst:                 filepath.Join(wd, "test-package"),
			generatePackageSpec: generatePackSpec,
			pathToExpectedFile:  filepath.Join(wd, "..", "..", "features", "test-package", "docker-compose.dev.yml"),
		},
	}

	for _, tc := range testCases {
		defer os.RemoveAll(tc.dst)

		err := os.Mkdir(tc.dst, os.ModePerm)
		jtest.RequireNil(t, err)

		err = createFileFromTemplate(tc.src, tc.dst, tc.generatePackageSpec)
		jtest.RequireNil(t, err)

		expectedData, err := os.ReadFile(tc.pathToExpectedFile)
		jtest.RequireNil(t, err)

		actualData, err := os.ReadFile(filepath.Join(tc.dst, tc.src))
		jtest.RequireNil(t, err)

		expected := crc32.ChecksumIEEE(expectedData)
		actual := crc32.ChecksumIEEE(actualData)

		require.Equal(t, expected, actual)

		// ensure removal after each test case
		err = os.RemoveAll(tc.dst)
		jtest.RequireNil(t, err)
	}
}

func TestGenerateConfigFile(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	type cases struct {
		config             core.Config
		pathToExpectedFile string
	}

	testCases := []cases{
		// case: assert config file created as expected
		{
			config: core.Config{
				ProjectName:   "test-project",
				Image:         "jembi/go-cli-test-image",
				PlatformImage: "jembi/platform:latest",
				LogPath:       "/tmp/logs",
				Packages:      []string{"client", "dashboard-visualiser-jsreport"},
				CustomPackages: []core.CustomPackage{
					{
						Id:   "disi-on-platform",
						Path: "git@github.com:jembi/disi-on-platform.git",
					},
				},
				Profiles: []core.Profile{
					{
						Name:     "env-var-test",
						Packages: []string{"dashboard-visualiser-jsreport", "disi-on-platform"},
						EnvVars:  []string{"SECOND=env_var_two_overwrite"},
						EnvFiles: []string{"../test-conf/.env.four"},
					},
				},
			},
			pathToExpectedFile: filepath.Join(wd, "..", "..", "features", "unit-test-configs", "config-case-6.yml"),
		},
		// case: assert invalid config file, missing field 'Image'
		{
			config: core.Config{
				ProjectName:   "test-project",
				PlatformImage: "jembi/platform:latest",
			},
		},
		// case: assert invalid config file, missing field 'ProjectName'
		{
			config: core.Config{
				Image:         "jembi/go-cli-test-image",
				PlatformImage: "jembi/platform:latest",
			},
		},
		// case: assert invalid config file, missing field 'PlatformImage'
		{
			config: core.Config{
				Image:       "jembi/go-cli-test-image",
				ProjectName: "test-project",
			},
		},
	}

	for _, tc := range testCases {
		defer os.Remove("config.yaml")

		err := GenerateConfigFile(&tc.config)
		if err != nil {
			require.Equal(t, ErrInvalidConfig.Error(), err.Error())
			continue
		}

		expectedConfigFile, err := os.Open(tc.pathToExpectedFile)
		jtest.RequireNil(t, err)

		expectedScanner := bufio.NewScanner(expectedConfigFile)
		expectedScanner.Split(bufio.ScanLines)

		var expected []byte
		for expectedScanner.Scan() {
			expected = append(expected, bytes.TrimSpace(expectedScanner.Bytes())...)
			expected = append(expected, '\n')
		}

		generatedConfigFile, err := os.Open("config.yaml")
		jtest.RequireNil(t, err)

		generatedScanner := bufio.NewScanner(generatedConfigFile)
		generatedScanner.Split(bufio.ScanLines)

		generated := []byte("---\n")
		for generatedScanner.Scan() {
			generated = append(generated, bytes.TrimSpace(generatedScanner.Bytes())...)
			generated = append(generated, []byte("\n")...)
		}

		expectedText := string(expected)
		generatedText := string(generated)

		require.Equal(t, expectedText, generatedText)
	}
}
