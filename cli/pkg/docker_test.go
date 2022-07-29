package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func Test_sliceContains(t *testing.T) {
	testCases := []struct {
		slice   []string
		element string
		result  bool
		name    string
	}{
		{
			name:    "SliceContain test - should return true when slice contains element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Optimus Prime",
			result:  true,
		},
		{
			name:    "SliceContain test - should return false when slice does not contain element",
			slice:   []string{"Optimus Prime", "Iron Hyde"},
			element: "Megatron",
			result:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ans := sliceContains(tt.slice, tt.element)

			if ans != tt.result {
				t.Fatal("SliceContains should return " + fmt.Sprintf("%t", tt.result) + " but returned " + fmt.Sprintf("%t", ans))
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_getPackagePaths(t *testing.T) {
	type args struct {
		inputArr []string
		flags    []string
	}
	tests := []struct {
		name             string
		args             args
		wantPackagePaths []string
	}{
		{
			name: "Test 1 - '-c' flag",
			args: args{
				inputArr: []string{"-c=../docs", "-c=./docs"},
				flags:    []string{"-c=", "--custom-package="},
			},
			wantPackagePaths: []string{"../docs", "./docs"},
		},
		{
			name: "Test 2 - '--custom-package' flag",
			args: args{
				inputArr: []string{"--custom-package=../docs", "--custom-package=./docs"},
				flags:    []string{"-c=", "--custom-package="},
			},
			wantPackagePaths: []string{"../docs", "./docs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPackagePaths := getPackagePaths(tt.args.inputArr, tt.args.flags); !assert.Equal(t, tt.wantPackagePaths, gotPackagePaths) {
				t.Fatalf("getPackagePaths() = %v, want %v", gotPackagePaths, tt.wantPackagePaths)
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_getEnvironmentVariables(t *testing.T) {
	type args struct {
		inputArr []string
		flags    []string
	}
	tests := []struct {
		name                     string
		args                     args
		wantEnvironmentVariables []string
	}{
		{
			name: "Test case environment variables found",
			args: args{
				inputArr: []string{"-e=NODE_ENV=PROD", "-e=DOMAIN_NAME=instant.com"},
				flags:    []string{"-e=", "--env-file="},
			},
			wantEnvironmentVariables: []string{"-e", "NODE_ENV=PROD", "-e", "DOMAIN_NAME=instant.com"},
		},
		{
			name: "Test case environment variables file found",
			args: args{
				inputArr: []string{"--env-file=../test.env"},
				flags:    []string{"-e=", "--env-file="},
			},
			wantEnvironmentVariables: []string{"--env-file", "../test.env"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEnvironmentVariables := getEnvironmentVariables(tt.args.inputArr, tt.args.flags); !assert.Equal(t, tt.wantEnvironmentVariables, gotEnvironmentVariables) {
				t.Fatalf("getEnvironmentVariables() = %v, want %v", gotEnvironmentVariables, tt.wantEnvironmentVariables)
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_extractCommands(t *testing.T) {
	customOptions.targetLauncher = "docker"

	testCases := []struct {
		startupCommands []string
		expectedResults CommandsOptions
		name            string
	}{
		{
			startupCommands: []string{"init", "-t=docker", "--image-version=v2.0.1", "-c=../test", "-c=../test1", "-e=NODE_ENV=dev", "-onlyFlag", "core"},
			expectedResults: CommandsOptions{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "init",
				otherFlags:           []string{"-onlyFlag"},
				targetLauncher:       "docker",
				packages:             []string{"core"},
				customPackagePaths:   []string{"../test", "../test1"},
				imageVersion:         "v2.0.1",
			},
			name: "Extract commands test 1 - should return the expected commands",
		},
		{
			startupCommands: []string{"up", "-t=kubernetes", "--image-version=v2.0.2", "-c=../test", "-c=../test1", "-e=NODE_ENV=dev", "-onlyFlag", "core"},
			expectedResults: CommandsOptions{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "up",
				otherFlags:           []string{"-onlyFlag"},
				targetLauncher:       "kubernetes",
				packages:             []string{"core"},
				customPackagePaths:   []string{"../test", "../test1"},
				imageVersion:         "v2.0.2",
			},
			name: "Extract commands test 2 - should return the expected commands",
		},
		{
			startupCommands: []string{"down", "-t=k8s", "--image-version=v2.0.2", "-c=../test", "-c=../test1", "--env-file=../test.env", "-onlyFlag", "core", "hapi-fhir"},
			expectedResults: CommandsOptions{
				environmentVariables: []string{"--env-file", "../test.env"},
				deployCommand:        "down",
				otherFlags:           []string{"-onlyFlag"},
				targetLauncher:       "k8s",
				packages:             []string{"core", "hapi-fhir"},
				customPackagePaths:   []string{"../test", "../test1"},
				imageVersion:         "v2.0.2",
			},
			name: "Extract commands test 3 - should return the expected commands",
		},
		{
			startupCommands: []string{"destroy", "-t=swarm", "--image-version=v2.0.2", "--custom-package=../test", "-c=../test1", "-e=NODE_ENV=dev", "--onlyFlag", "core", "hapi-fhir"},
			expectedResults: CommandsOptions{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "destroy",
				otherFlags:           []string{"--onlyFlag"},
				targetLauncher:       "swarm",
				packages:             []string{"core", "hapi-fhir"},
				customPackagePaths:   []string{"../test", "../test1"},
				imageVersion:         "v2.0.2",
			},
			name: "Extract commands test 4 - should return the expected commands",
		},
		{
			startupCommands: []string{"destroy", "--image-version=v2.0.2", "--custom-package=../test", "-c=../test1", "-e=NODE_ENV=dev", "--onlyFlag", "core", "hapi-fhir"},
			expectedResults: CommandsOptions{
				environmentVariables: []string{"-e", "NODE_ENV=dev"},
				deployCommand:        "destroy",
				otherFlags:           []string{"--onlyFlag"},
				targetLauncher:       "docker",
				packages:             []string{"core", "hapi-fhir"},
				customPackagePaths:   []string{"../test", "../test1"},
				imageVersion:         "v2.0.2",
			},
			name: "Extract commands test 4 - should return the expected commands",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			commandOptions := extractCommands(tt.startupCommands)

			if !assert.Equal(t, tt.expectedResults.environmentVariables, commandOptions.environmentVariables) {
				t.Fatal("ExtractCommands should return the correct environment variables")
			}
			if !assert.Equal(t, tt.expectedResults.deployCommand, commandOptions.deployCommand) {
				t.Fatal("ExtractCommands should return the correct deploy command")
			}
			if !assert.Equal(t, tt.expectedResults.otherFlags, commandOptions.otherFlags) {
				t.Fatal("ExtractCommands should return the correct 'otherFlags'")
			}
			if !assert.Equal(t, tt.expectedResults.targetLauncher, commandOptions.targetLauncher) {
				t.Fatal("ExtractCommands should return the correct targetLauncher")
			}
			if !assert.Equal(t, tt.expectedResults.packages, commandOptions.packages) {
				t.Fatal("ExtractCommands should return the correct packages")
			}
			if !assert.Equal(t, tt.expectedResults.customPackagePaths, commandOptions.customPackagePaths) {
				t.Fatal("ExtractCommands should return the correct custom package paths")
			}
			if !assert.Equal(t, tt.expectedResults.imageVersion, commandOptions.imageVersion) {
				t.Fatal("ExtractCommands should return the correct image version")
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_createZipFile(t *testing.T) {
	var reader io.Reader

	type args struct {
		file    string
		content io.Reader
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockOsCreate func(name string) (*os.File, error)
		mockIoCopy   func(dst io.Writer, src io.Reader) (written int64, err error)
	}{
		{
			name: "Test case create zip file no errors",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: false,
			mockOsCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, nil
			},
		},
		{
			name: "Test case create zip file with errors from OsCreate",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			mockOsCreate: func(name string) (*os.File, error) {
				return &os.File{}, errors.New("Test error")
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, nil
			},
		},
		{
			name: "Test case create zip file with errors from IoCopy",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			mockOsCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 1, errors.New("Test error")
			},
		},
		{
			name: "Test case create empty zip",
			args: args{
				file:    "test_zip.zip",
				content: reader,
			},
			wantErr: true,
			mockOsCreate: func(name string) (*os.File, error) {
				return &os.File{}, nil
			},
			mockIoCopy: func(dst io.Writer, src io.Reader) (written int64, err error) {
				return 0, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OsCreate = tt.mockOsCreate
			IoCopy = tt.mockIoCopy

			if err := createZipFile(tt.args.file, tt.args.content); (err != nil) != tt.wantErr {
				t.Fatalf("createZipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func Test_runCommand(t *testing.T) {
	testCases := []struct {
		commandName     string
		suppressErrors  []string
		commandSlice    []string
		pathToPackage   string
		errorString     error
		name            string
		mockExecCommand func(commandName string, commandSlice ...string) *exec.Cmd
	}{
		{
			commandName:     "docker",
			suppressErrors:  nil,
			commandSlice:    []string{"ps"},
			pathToPackage:   "",
			errorString:     nil,
			name:            "runCommand - run basic docker ps test",
			mockExecCommand: exec.Command,
		},
		{
			commandName:     "docker",
			suppressErrors:  nil,
			commandSlice:    []string{"volume", "rm", "test-volume"},
			pathToPackage:   "",
			errorString:     fmt.Errorf("Error waiting for Cmd. Error: No such volume: test-volume: exit status 1"),
			name:            "runCommand - removing nonexistant volume should return error",
			mockExecCommand: exec.Command,
		},
		{
			commandName:     "docker",
			suppressErrors:  []string{"Error: No such volume: test-volume"},
			commandSlice:    []string{"volume", "rm", "test-volume"},
			pathToPackage:   "",
			errorString:     nil,
			name:            "runCommand - error thrown should be suppressed",
			mockExecCommand: exec.Command,
		},
		{
			commandName:    "git",
			suppressErrors: nil,
			commandSlice:   []string{"clone", "git@github.com:testhie/test.git"},
			pathToPackage:  "test",
			errorString:    nil,
			name:           "runCommand - clone a custom package and return its location",
			mockExecCommand: func(commandName string, commandSlice ...string) *exec.Cmd {
				cmd := exec.Command("pwd")
				return cmd
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = tt.mockExecCommand
			pathToPackage, err := runCommand(tt.commandName, tt.suppressErrors, tt.commandSlice...)
			if !assert.Equal(t, tt.pathToPackage, pathToPackage) {
				t.Fatal("RunCommand failed - path to package returned is incorrect " + pathToPackage)
			}
			if err != nil && tt.errorString != nil && !assert.Equal(t, tt.errorString.Error(), err.Error()) {
				t.Fatal("RunCommand failed - error returned incorrect")
			}

			if (err != nil && tt.errorString == nil) || (err == nil && tt.errorString != nil) {
				t.Log("Expected:", tt.errorString)
				t.Log("Actual:", err.Error())
				t.Fatal("RunCommand failed - error returned incorrect")
			}

			t.Log(tt.name + " passed!")
		})
	}
}

func Test_mountPackage(t *testing.T) {
	defer gock.Off()
	defer resetMountPackageMocks(runCommand, unzipPackage, untarPackage)

	testCases := []struct {
		pathToPackage    string
		mockServer       func()
		name             string
		wantErr          bool
		errorString      string
		mockRunCommand   func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error)
		mockUnzipPackage func(zipContent io.ReadCloser) (pathToPackage string, err error)
		mockUntarPackage func(tarContent io.ReadCloser) (pathToPackage string, err error)
	}{
		{
			pathToPackage: "http://test:8080/test",
			mockServer:    func() {},
			errorString:   "Error in downloading custom package",
			wantErr:       true,
			name:          "mountPackage - should return error when downloading custom package fails",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "git@github.com:test/test.git",
			mockServer:    func() {},
			errorString:   "Error in git cloning",
			wantErr:       true,
			name:          "mountPackage - should return error when 'git cloning' a custom package fails",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", errors.New("Error in git cloning")
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "http://test.com/test.zip",
			mockServer: func() {
				gock.New("http://test.com").
					Get("/test.zip").
					Reply(200).
					BodyString("Zip File Content")
			},
			errorString: "Error in unzipping package",
			wantErr:     true,
			name:        "mountPackage - should return error when unziping the custom package fails",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", errors.New("Error in unzipping package")
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "http://test.com/test.tar.gz",
			mockServer: func() {
				gock.New("http://test.com").
					Get("/test.tar.gz").
					Reply(200).
					BodyString("Tar File Content")
			},
			errorString: "Error in untarring package",
			wantErr:     true,
			name:        "mountPackage - should return error when untarring the custom package fails",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", errors.New("Error in untarring package")
			},
		},
		{
			pathToPackage: "http://test.com/test.tar.gz",
			mockServer: func() {
				gock.New("http://test.com").
					Get("/test.tar.gz").
					Reply(200).
					BodyString("Tar File Content")
			},
			errorString: "Error in copying package",
			wantErr:     true,
			name:        "mountPackage - should return error when copying the custom package to the instant docker container fails",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", errors.New("Error in copying package")
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "git@github.com:test/test.git",
			mockServer:    func() {},
			errorString:   "",
			wantErr:       false,
			name:          "mountPackage - should git clone custom package and copy it to the instant docker container",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				pathToPackage = "./test"
				if commandName == "git" {
					return pathToPackage, nil
				}
				if commandName == "docker" && commandSlice[1] != pathToPackage {
					t.Fatal("Path to package returned is incorrect")
				}
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "http://test.com/test1.zip",
			mockServer: func() {
				gock.New("http://test.com").
					Get("/test1.zip").
					Reply(200).
					BodyString("Zip File Content")
			},
			errorString: "",
			wantErr:     false,
			name:        "mountPackage - should unzip custom package and copy it to the instant docker container",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				pathToPackage = "./test1"

				if commandName == "docker" && commandSlice[1] != pathToPackage {
					t.Fatal("Path to custom package returned is incorrect")
				}
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "./test1", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
		},
		{
			pathToPackage: "http://test.com/test2.tar.gz",
			mockServer: func() {
				gock.New("http://test.com").
					Get("/test2.tar.gz").
					Reply(200).
					BodyString("Tar File Content")
			},
			errorString: "",
			wantErr:     false,
			name:        "mountPackage - should untar custom package and copy it to the instant docker container",
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				pathToPackage = "./test2"

				if commandName == "docker" && commandSlice[1] != pathToPackage {
					t.Fatal("Path to custom package returned is incorrect")
				}
				return "", nil
			},
			mockUnzipPackage: func(zipContent io.ReadCloser) (pathToPackage string, err error) {
				return "", nil
			},
			mockUntarPackage: func(tarContent io.ReadCloser) (pathToPackage string, err error) {
				return "./test2", nil
			},
		},
	}

	for _, tt := range testCases {
		runCommand = tt.mockRunCommand
		tt.mockServer()
		unzipPackage = tt.mockUnzipPackage
		untarPackage = tt.mockUntarPackage

		t.Run(tt.name, func(t *testing.T) {
			err := mountCustomPackage(tt.pathToPackage)
			if err == nil && tt.wantErr {
				t.Fatal("Expected error - '" + tt.errorString + "' but got nil")
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errorString) {
				t.Fatal(err.Error())
			}
			t.Log(tt.name + " passed!")
		})
	}
}

func resetMountPackageMocks(originalRunCommand func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error), originalUnzip func(zipContent io.ReadCloser) (pathToPackage string, err error), originalUntar func(tarContent io.ReadCloser) (pathToPackage string, err error)) {
	runCommand = originalRunCommand
	unzipPackage = originalUnzip
	untarPackage = originalUntar
}

func Test_unzipPackage(t *testing.T) {
	zipFile := createTestZipFile()
	defer zipFile.Close()

	type args struct {
		zipContent io.ReadCloser
	}
	tests := []struct {
		name              string
		args              args
		wantPathToPackage string
		wantErr           bool
		mockZipOpenReader func(name string) (*zip.ReadCloser, error)
		mockOsMkdirAll    func(path string, perm fs.FileMode) error
		mockFilepathJoin  func(elem ...string) string
		mockOsOpenFile    func(name string, flag int, perm fs.FileMode) (*os.File, error)
		mockOsRemove      func(name string) error
	}{
		{
			name: "Test case no errors",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "test_zip.zip",
			wantErr:           false,
			mockZipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zipReader, err := zip.OpenReader(zipFile.Name())
				if err != nil {
					t.Fatal(err)
				}

				file := new(zip.File)
				file.Name = "./"
				zipReader.File = append(zipReader.File, file)
				return zipReader, nil
			},
			mockFilepathJoin: func(elem ...string) string {
				return zipFile.Name()
			},
			mockOsMkdirAll: func(path string, perm fs.FileMode) error {
				return nil
			},
			mockOsOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return &os.File{}, nil
			},
			mockOsRemove: func(name string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from ZipOpenReader",
			args: args{
				zipContent: zipFile,
			},
			wantPathToPackage: "",
			wantErr:           true,
			mockZipOpenReader: func(name string) (*zip.ReadCloser, error) {
				return &zip.ReadCloser{}, errors.New("ZipOpenReader error")
			},
		},
		{
			name: "Test case receive error from OsMkdirAll",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			mockZipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zipReader := new(zip.ReadCloser)
				defer zipReader.Close()

				file := new(zip.File)
				file.Name = "./"
				zipReader.File = append(zipReader.File, file)
				return zipReader, nil
			},
			mockFilepathJoin: func(elem ...string) string {
				return ""
			},
			mockOsMkdirAll: func(path string, perm fs.FileMode) error {
				return errors.New("OsMkdirAll error")
			},
		},
		{
			name: "Test case receive error from OsOpenFile",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			mockZipOpenReader: func(name string) (*zip.ReadCloser, error) {
				zipReader, err := zip.OpenReader(zipFile.Name())
				if err != nil {
					t.Fatal(err)
				}
				return zipReader, nil
			},
			mockFilepathJoin: func(elem ...string) string {
				return ""
			},
			mockOsMkdirAll: func(path string, perm fs.FileMode) error {
				return nil
			},
			mockOsOpenFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return &os.File{}, errors.New("OsOpenFile error")
			},
		},
		{
			name: "Test case receive error from OsRemove",
			args: args{
				zipContent: &io.PipeReader{},
			},
			wantPathToPackage: "",
			wantErr:           true,
			mockZipOpenReader: func(name string) (*zip.ReadCloser, error) {
				return &zip.ReadCloser{}, nil
			},
			mockFilepathJoin: func(elem ...string) string {
				return ""
			},
			mockOsRemove: func(name string) error {
				return errors.New("OsRemove error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipFileExists(zipFile.Name())
			defer cleanFiles(zipFile.Name())
			defaultMockVariables()

			ZipOpenReader = tt.mockZipOpenReader
			FilepathJoin = tt.mockFilepathJoin
			OsMkdirAll = tt.mockOsMkdirAll
			OsOpenFile = tt.mockOsOpenFile
			OsRemove = tt.mockOsRemove

			gotPathToPackage, err := unzipPackage(tt.args.zipContent)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unzipPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotPathToPackage != tt.wantPathToPackage {
				t.Fatalf("unzipPackage() = %v, want %v", gotPathToPackage, tt.wantPathToPackage)
			}
			t.Log(tt.name, " passed!")
		})
	}
}

func createTestZipFile() *os.File {
	zipFile, err := os.Create("test_zip.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	writer, err := zipWriter.Create("test_txt.txt")
	if err != nil {
		panic(err)
	}

	_, err = writer.Write([]byte("test data"))
	if err != nil {
		panic(err)
	}

	err = zipWriter.Flush()
	if err != nil {
		panic(err)
	}

	return zipFile
}

func zipFileExists(fileName string) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		createTestZipFile()
	}
}

func cleanFiles(fileName string) {
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		os.Remove(fileName)
	}
}

func defaultMockVariables() {
	OsCreate = func(name string) (*os.File, error) {
		return &os.File{}, nil
	}
	IoCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
		return 1, nil
	}
}

func TestRunDeployCommand(t *testing.T) {
	type args struct {
		startupCommands []string
	}
	tests := []struct {
		name                   string
		args                   args
		wantErr                bool
		mockRunCommand         func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error)
		mockMountCustomPackage func(pathToPackage string) error
	}{
		{
			name: "Test case expect no errors",
			args: args{
				startupCommands: []string{"init", "core", "-c=./local/cPack", "--image-version=latest", "-t=docker"},
			},
			wantErr: false,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", nil
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from first call to RunCommand()",
			args: args{
				startupCommands: []string{"init", "core", "-c=./local/cPack", "--image-version=latest", "-t=docker"},
			},
			wantErr: true,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", errors.New("test error")
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from second call to RunCommand()",
			args: args{
				startupCommands: []string{"down", "--image-version=latest", "-t=docker"},
			},
			wantErr: true,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", errors.New("test error")
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return nil
			},
		},
		{
			name: "Test case ignore error from third call to RunCommand()",
			args: args{
				startupCommands: []string{"down", "--image-version=latest", "-t=docker"},
			},
			wantErr: false,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				if commandSlice[0] == "start" {
					return "", errors.New("test error")
				}
				return "", nil
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return nil
			},
		},
		{
			name: "Test case verify commandSlice append",
			args: args{
				startupCommands: []string{"up", "hmis", "mcsd", "--env-file=./home/bin", "-e=NODE_ENV=DEV",
					"-e=DOMAIN_NAME=instant.com", "-c=./usr/local/cPack", "--only", "--dev", "--image-version=v1.03a", "-t=k8s"},
			},
			wantErr: false,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				if commandSlice[0] != "create" {
					return "", nil
				}

				expectedCommandSlice := []string{
					"create",
					"--rm",
					"--mount=type=volume,src=instant,dst=/instant",
					"--name", "instant-openhie",
					"-v", "/var/run/docker.sock:/var/run/docker.sock",
					"--network", "host",
					"--env-file", "./home/bin", "-e", "NODE_ENV=DEV", "-e", "DOMAIN_NAME=instant.com",
					":v1.03a", "up", "--only", "--dev",
					"-t", "k8s", "hmis", "mcsd",
				}

				if !assert.Equal(t, expectedCommandSlice, commandSlice) {
					t.Fatal("commandSlice not matched")
				}

				return "", nil
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return nil
			},
		},
		{
			name: "Test case receive error from MountCustomPackage()",
			args: args{
				startupCommands: []string{"init", "core", "-c=./local/cPack", "--image-version=latest", "-t=docker"},
			},
			wantErr: true,
			mockRunCommand: func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
				return "", nil
			},
			mockMountCustomPackage: func(pathToPackage string) error {
				return errors.New("test error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunCommand = tt.mockRunCommand
			MountCustomPackage = tt.mockMountCustomPackage

			if err := RunDeployCommand(tt.args.startupCommands); (err != nil) != tt.wantErr {
				t.Fatalf("RunDeployCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(tt.name, " passed!")
		})
	}
}
