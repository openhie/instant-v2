package generate

import (
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

		require.Equal(t, crc32.ChecksumIEEE(expectedData), crc32.ChecksumIEEE(actualData))

		// ensure removal after each test case
		err = os.RemoveAll(tc.dst)
		jtest.RequireNil(t, err)
	}
}
