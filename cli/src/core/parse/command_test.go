package parse

import (
	"cli/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCustomPackageName(t *testing.T) {
	type cases struct {
		customPackage core.CustomPackage
		expectName    string
	}

	testCases := []cases{
		// case: remote custom package with ID
		{core.CustomPackage{Id: "test", Path: "https://github.com/test/test-package.git"}, "test"},

		// case: local custom package with ID
		{core.CustomPackage{Id: "test-package", Path: "../../home/test-package"}, "test-package"},

		// case: git custom package
		{core.CustomPackage{Path: "https://github.com/test/test-package.git"}, "test-package"},

		// case: tar custom package
		{core.CustomPackage{Path: "https://github.com/test/test-package.tar"}, "test-package"},

		// case: zip custom package
		{core.CustomPackage{Path: "https://github.com/test/test-package.zip"}, "test-package"},

		// case: custom package without ID
		{core.CustomPackage{Path: "/home/path/test-package"}, "test-package"},

		// case: git custom package without full URL
		{core.CustomPackage{Path: "git@github.com:test/test-package.git"}, "test-package"},
	}

	for _, tc := range testCases {
		got := GetCustomPackageName(tc.customPackage)
		require.Equal(t, tc.expectName, got)
	}
}

func TestGetInstantCommand(t *testing.T) {
	packageOnlyPackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
	}
	upPackageSpec := core.PackageSpec{
		DeployCommand: "up",
		Packages:      []string{"test-package"},
	}
	downPackageSpec := core.PackageSpec{
		DeployCommand: "down",
		Packages:      []string{"test-package"},
	}
	destroyPackageSpec := core.PackageSpec{
		DeployCommand: "destroy",
		Packages:      []string{"test-package"},
	}
	packagesOnlyPackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package", "another-test-package"},
	}
	devFlagTruePackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsDev:         true,
	}
	devFlagFalsePackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
	}
	onlyFlagTruePackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsOnly:        true,
	}
	onlyFlagFalsePackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
	}
	customPackageOnlyPackageSpec := core.PackageSpec{
		DeployCommand: "init",
		CustomPackages: []core.CustomPackage{
			{
				Id:   "custom-package",
				Path: "../custom-package",
			},
		},
	}
	packageWithCustomPackagePackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		CustomPackages: []core.CustomPackage{
			{
				Id:   "custom-package",
				Path: "../custom-package",
			},
		},
	}
	fullPackageSpec := core.PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		CustomPackages: []core.CustomPackage{
			{
				Id:   "custom-package",
				Path: "../custom-package",
			},
		},
		IsDev:  true,
		IsOnly: true,
	}

	tables := []struct {
		description    string
		packageSpec    core.PackageSpec
		instantCommand []string
	}{
		{"init package only", packageOnlyPackageSpec, []string{"init", "-t", "swarm", "test-package"}},
		{"up package only", upPackageSpec, []string{"up", "-t", "swarm", "test-package"}},
		{"down package only", downPackageSpec, []string{"down", "-t", "swarm", "test-package"}},
		{"destroy package only", destroyPackageSpec, []string{"destroy", "-t", "swarm", "test-package"}},
		{"init packages only", packagesOnlyPackageSpec, []string{"init", "-t", "swarm", "test-package", "another-test-package"}},
		{"init dev flag true", devFlagTruePackageSpec, []string{"init", "-t", "swarm", "--dev", "test-package"}},
		{"init dev flag false", devFlagFalsePackageSpec, []string{"init", "-t", "swarm", "test-package"}},
		{"init only flag true", onlyFlagTruePackageSpec, []string{"init", "-t", "swarm", "--only", "test-package"}},
		{"init only flag false", onlyFlagFalsePackageSpec, []string{"init", "-t", "swarm", "test-package"}},
		{"init custom package only", customPackageOnlyPackageSpec, []string{"init", "-t", "swarm", "custom-package"}},
		{"init package with custom package", packageWithCustomPackagePackageSpec, []string{"init", "-t", "swarm", "test-package", "custom-package"}},
		{"init package and custom package with dev and only flag", fullPackageSpec, []string{"init", "-t", "swarm", "--dev", "--only", "test-package", "custom-package"}},
	}

	for _, table := range tables {
		instantCommand := GetInstantCommand(table.packageSpec)
		require.Equal(t, instantCommand, table.instantCommand)
	}
}
