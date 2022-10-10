package core

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"testing"

	"github.com/docker/docker/api/types"
	_container "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	}
	onlyFlagTruePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
		IsOnly:        true,
	}
	onlyFlagFalsePackageSpec := PackageSpec{
		DeployCommand: "init",
		Packages:      []string{"test-package"},
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
		assert.Equal(t, instantCommand, table.instantCommand)
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
		assert.Equal(t, customPackageName, table.customPackageName)
	}
}

func Test_attachUntilRemoved(t *testing.T) {
	mockApiClient := new(MockApiClient)

	// Case: receive error from cli.Container attach
	mockApiClient.On("ContainerAttach").Return("test error").Once()
	err := attachUntilRemoved(mockApiClient, context.Background(), "")
	jtest.Require(t, errors.New("test error"), err)

	// Case: receive success message to successChannel
	mockApiClient.On("ContainerAttach").Return(nil)

	mockApiClient.On("ContainerWait").Return(_container.ContainerWaitOKBody{
		StatusCode: 0,
		Error:      nil,
	}, nil).Once()

	err = attachUntilRemoved(mockApiClient, context.Background(), "")
	jtest.RequireNil(t, err)

	// Case: receive expected "No such container" message
	mockApiClient.On("ContainerWait").Return(nil, "No such container").Once()

	err = attachUntilRemoved(mockApiClient, context.Background(), "")
	jtest.RequireNil(t, err)

	// Case: receive error to errorChannel
	mockApiClient.On("ContainerWait").Return(nil, "test error").Once()

	err = attachUntilRemoved(mockApiClient, context.Background(), "")
	jtest.Require(t, errors.New("test error"), err)
}

type MockApiClient struct {
	mock.Mock
	client.ContainerAPIClient
}

func (mock *MockApiClient) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	args := mock.Called()

	response := types.HijackedResponse{
		Conn:   &net.IPConn{},
		Reader: bufio.NewReader(bytes.NewReader([]byte("test"))),
	}

	var err error
	if args.Get(0) != nil {
		err = errors.New(args.Get(0).(string))
	}

	return response, err
}

func (mock *MockApiClient) ContainerWait(ctx context.Context, container string, condition _container.WaitCondition) (<-chan _container.ContainerWaitOKBody, <-chan error) {
	args := mock.Called()

	newChan := make(chan _container.ContainerWaitOKBody)

	var containerWaitBody _container.ContainerWaitOKBody
	b := args.Get(0)
	if b != nil {
		containerWaitBody = args.Get(0).(_container.ContainerWaitOKBody)

		go func() {
			newChan <- containerWaitBody
		}()
	}

	errChan := make(chan error)

	var err error
	c := args.Get(1)
	if c != nil {
		err = errors.New(args.Get(1).(string))

		go func() {
			errChan <- err
		}()
	}

	return newChan, errChan
}

func Test_getCustomPackageName(t *testing.T) {
	type cases struct {
		customPackage CustomPackage
		expectName    string
	}

	testCases := []cases{
		{CustomPackage{Path: "https://github.com/jembi/instant-openhie-template-package.git"}, "instant-openhie-template-package"},
		{CustomPackage{Path: "https://github.com/jembi/instant-openhie-template-package.tar"}, "instant-openhie-template-package"},
		{CustomPackage{Path: "https://github.com/jembi/instant-openhie-template-package.zip"}, "instant-openhie-template-package"},
		{CustomPackage{Id: "template", Path: "https://github.com/jembi/instant-openhie-template-package.git"}, "template"},
	}

	for _, tc := range testCases {
		got := getCustomPackageName(tc.customPackage)
		assert.Equal(t, tc.expectName, got)
	}
}
