package core

import "testing"

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestGetInstantCommand(t *testing.T) {
	packageOnlyPackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
	}
	upPackageSpec := PackageSpec{
		DeployCommand: "up",
		Packages:      []string{"test-package"},
	}
	downPackageSpec := PackageSpec{
		DeployCommand: "down",
		Packages:      []string{"test-package"},
	}
	destroyPackageSpec := PackageSpec{
		DeployCommand: "destroy",
		Packages:      []string{"test-package"},
	}
	packagesOnlyPackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package", "another-test-package"},
	}
	devFlagTruePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsDev:         true,
	}
	devFlagFalsePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsDev:         false,
	}
	onlyFlagTruePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsOnly:        true,
	}
	onlyFlagFalsePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsOnly:        false,
	}
	customPackageOnlyPackageSpec := PackageSpec{
		DeployCommand: "init",
		CustomPackages: []CustomPackage{
			{
				Id:   "custom-package",
				Path: "../custom-package",
			},
		},
	}
	packageWithCustomPackagePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		CustomPackages: []CustomPackage{
			{
				Id:   "custom-package",
				Path: "../custom-package",
			},
		},
	}
	fullPackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		CustomPackages: []CustomPackage{
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
		packageSpec    PackageSpec
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
		instantCommand := getInstantCommand(table.packageSpec)
		if !equal(instantCommand, table.instantCommand) {
			t.Errorf("Failed: %v, expected: %v , received: %v", table.description, table.instantCommand, instantCommand)
		}
	}
}

func TestGetCustomPackageName(t *testing.T) {
	idCustomPackage := CustomPackage{
		Id:   "test-package",
		Path: "../test-package",
	}
	pathAbsoluteCustomPackage := CustomPackage{
		Path: "/home/path/test-package",
	}
	pathRelativeCustomPackage := CustomPackage{
		Path: "../test-package",
	}
	pathGitCustomPackage := CustomPackage{
		Path: "git@github.com:test/test-package.git",
	}
	pathZipCustomPackage := CustomPackage{
		Path: "https://github.com/test/test-package.zip",
	}
	pathTarCustomPackage := CustomPackage{
		Path: "https://github.com/test/test-package.tar",
	}
	tables := []struct {
		description       string
		customPackage     CustomPackage
		customPackageName string
	}{
		{"custom package with id", idCustomPackage, "test-package"},
		{"custom package with absolute path", pathAbsoluteCustomPackage, "test-package"},
		{"custom package with relative path", pathRelativeCustomPackage, "test-package"},
		{"custom package with git path", pathGitCustomPackage, "test-package"},
		{"custom package with http zip path", pathZipCustomPackage, "test-package"},
		{"custom package with http tar path", pathTarCustomPackage, "test-package"},
	}

	for _, table := range tables {
		customPackageName := getCustomPackageName(table.customPackage)
		if customPackageName != table.customPackageName {
			t.Errorf("Failed: %v, expected: %v , received: %v", table.description, table.customPackageName, customPackageName)
		}
	}
}
