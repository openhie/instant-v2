package pkg

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"

	viperUtil "cli/cmd/util"
	"cli/core"

	"github.com/docker/docker/api/types"
	"github.com/luno/jettison/jtest"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var configCaseOne = core.Config{
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
}

func Test_unmarshalConfig(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	// case: match configCaseOne
	configViper, err := viperUtil.SetConfigViper(wd + "/../../features/unit-test-configs/config-case-1.yml")
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	assert.Equal(t, configCaseOne, *config)

	// case: return invalid config file syntax error
	configViper, err = viperUtil.SetConfigViper(wd + "/../../features/unit-test-configs/config-case-6.yml")
	jtest.RequireNil(t, err)

	_, err = unmarshalConfig(core.Config{}, configViper)
	assert.Equal(t, ErrInvalidConfigFileSyntax.Error(), err.Error())
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
		// case: return error from non-existant env file directory
		{
			profileName:         "bad-env-file-path",
			configFilePath:      wd + "/../../features/unit-test-configs/config-case-3.yml",
			expectedErrorString: "stat ./features/test-conf/.env.tests: no such file or directory",
		},
		// case: no profile specified, dev flag specified, return nil error
		{
			boolFlagName:   "dev",
			configFilePath: wd + "/../../features/test-conf/config.yml",
		},
		// case: no profile specified, only flag specified, return nil error
		{
			boolFlagName:   "only",
			configFilePath: wd + "/../../features/test-conf/config.yml",
		},
	}

	for _, tc := range testCases {
		cmd, config := setupLoadInProfileParams(t, tc.configFilePath)

		setPackageActionFlags(cmd)
		if tc.boolFlagName != "" {
			setupBoolFlags(t, cmd, tc.boolFlagName)
		}

		if tc.profileName != "" {
			cmd.Flags().Set("profile", tc.profileName)
		}

		_, err = loadInProfileParams(cmd, *config, core.PackageSpec{})
		if err != nil && !assert.Equal(t, tc.expectedErrorString, err.Error()) {
			t.FailNow()
		} else if tc.expectedErrorString == "" && err != nil {
			t.FailNow()
		}
	}

	// case: load in environment variables from more than one env file
	cmd, config := setupLoadInProfileParams(t, wd+"/../../features/unit-test-configs/config-case-2.yml")
	setPackageActionFlags(cmd)
	cmd.Flags().Set("profile", "non-only")

	packageSpec, err := loadInProfileParams(cmd, *config, core.PackageSpec{})
	jtest.RequireNil(t, err)

	sort.Slice(packageSpec.EnvironmentVariables, func(i, j int) bool {
		return strings.Contains(packageSpec.EnvironmentVariables[i], "FIRST_ENV_VAR")
	})

	assert.Equal(t, packageSpec.EnvironmentVariables, []string{"FIRST_ENV_VAR=number_one", "SECOND_ENV_VAR=number_two"})
}

func setupLoadInProfileParams(t *testing.T, configFilePath string) (*cobra.Command, *core.Config) {
	configViper, err := viperUtil.SetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	return &cobra.Command{}, config
}

func setupBoolFlags(t *testing.T, cmd *cobra.Command, boolFlagName string) {
	err := cmd.Flags().Set(boolFlagName, "true")
	jtest.RequireNil(t, err)
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

	configViper, err := viperUtil.SetConfigViper(wd + "/../../features/unit-test-configs/config-case-4.yml")
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	gotCustomPackages := getCustomPackages(config, []string{"path-to-1", "path-to-2"})

	assert.Equal(t, expectedCustomPackages, gotCustomPackages)
}

func Test_getPackageSpecFromParams(t *testing.T) {
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
			configFilePath: wd + "/../../features/unit-test-configs/config-case-2.yml",
			hookFunc: func(cmd *cobra.Command) {
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
			configFilePath: wd + "/../../features/unit-test-configs/config-case-2.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().Set("name", "pack-1")
				cmd.Flags().Set("env-file", wd+"/../../features/test-conf/awlikdeuh")
			},
			errorString: "no such file or directory",
		},
		// case: return no error when not specifying an env-file
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-2.yml",
			hookFunc: func(cmd *cobra.Command) {
				cmd.Flags().Set("name", "pack-1")
			},
		},
		// case: place .env file in main dir, but don't use its env vars
		{
			configFilePath: wd + "/../../features/unit-test-configs/config-case-2.yml",
			hookFunc: func(cmd *cobra.Command) {
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

		pSpec, err := getPackageSpecFromParams(cmd, config)
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
	configViper, err := viperUtil.SetConfigViper(configFilePath)
	jtest.RequireNil(t, err)

	config, err := unmarshalConfig(core.Config{}, configViper)
	jtest.RequireNil(t, err)

	cmd := &cobra.Command{}
	setPackageActionFlags(cmd)

	cmd.Flags().StringSlice("env-file", []string{""}, "")

	hookFunc(cmd)

	return cmd, config
}

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
		if !assert.Equal(t, tc.wantImageName, tc.config.Image) {
			t.FailNow()
		}
	}
}

func Test_hasImage(t *testing.T) {
	type cases struct {
		imageName string

		images    []types.ImageSummary
		wantMatch bool
	}

	testCases := []cases{
		// case: no match
		{
			imageName: "matchImage",
			images: []types.ImageSummary{
				{
					RepoTags: []string{"no-match-1", "no-match-2"},
				},
			},
			wantMatch: false,
		},
		// case: match
		{
			imageName: "matchImage",
			images: []types.ImageSummary{
				{
					RepoTags: []string{"no-match-1", "no-match-2", "matchImage"},
				},
			},
			wantMatch: true,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.wantMatch, hasImage(tc.imageName, tc.images)) {
			t.FailNow()
		}
	}
}

func copyFile(src, dst string) error {
	file, err := os.OpenFile(src, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	copiedFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(copiedFile)

	for fileScanner.Scan() {
		b := []byte(fileScanner.Text())

		_, err = writer.Write(append(b, []byte("\n")...))
		if err != nil {
			return err
		}
	}
	file.Close()

	err = writer.Flush()
	if err != nil {
		return err
	}

	return copiedFile.Close()
}
