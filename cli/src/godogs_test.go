package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
)

var (
	binaryFilePath string
	logs           string
	directoryNames []string
	customPackages = make(map[string]bool)
)

func theCommandIsRun(command string) error {
	customPackageNames, ok := hasCustomPackage(command)
	if ok {
		for _, customPackageName := range customPackageNames {
			// TODO: find a way to make the below code OS-agnostic (/tmp is linux specific)
			go monitorDirFor(filepath.Join("/tmp", "custom-package"), customPackageName)
		}
	}

	res, err := runTestCommand(binaryFilePath, strings.Split(command, " ")...)
	if err == nil {
		logs = res
	}
	return err
}

func theCommandIsRunWithProfile(command string, packages *godog.Table) error {
	if len(packages.Rows) > 0 {
		for i := 1; i < len(packages.Rows); i++ {
			// TODO: find a way to make the below code OS-agnostic (/tmp is linux specific)
			go monitorDirFor(filepath.Join("/tmp", "custom-package"), packages.Rows[i].Cells[0].Value)
		}
	}

	res, err := runTestCommand(binaryFilePath, strings.Split(command, " ")...)
	if err != nil {
		return err
	}
	logs = res

	return nil
}

func checkTheCLIOutputIs(expectedOutput string) error {
	return compareLogsAndOutputs(logs, expectedOutput)
}

func checkCustomPackages(packages *godog.Table) error {
	head := packages.Rows[0].Cells

	for i := 1; i < len(packages.Rows); i++ {
		for n, cell := range packages.Rows[i].Cells {
			switch head[n].Value {
			case "directory":
				directoryNames = append(directoryNames, cell.Value)

				if customPackages[cell.Value] {
					return nil
				}
			default:
				return errors.New("Unexpected column name: " + head[n].Value)
			}
		}
	}

	var v string

	for i := 1; i < len(packages.Rows); i++ {
		for _, cell := range packages.Rows[i].Cells {
			if !customPackages[cell.Value] {
				v += cell.Value + "\n"
			}
		}
	}

	return errors.New("did not create custom packages:\n" + v)
}

func compareLogsAndOutputs(inputLogs, expectedOutput string) error {
	if !strings.Contains(inputLogs, expectedOutput) {
		return errors.New("Logs received: '" + inputLogs + "\nSubstring expected: " + expectedOutput)
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	defer cleanUp()

	suite := &godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			binaryFilePath = buildBinary()
			copyFiles()

			sc.Step(`^check the CLI output is "([^"]*)"$`, checkTheCLIOutputIs)
			sc.Step(`^the command "([^"]*)" is run$`, theCommandIsRun)
			sc.Step(`^the command "([^"]*)" is run with profile$`, theCommandIsRunWithProfile)
			sc.Step(`^check that the CLI added custom packages$`, checkCustomPackages)
		},
	}

	if suite.Run() != 0 {
		fmt.Println("Tests failed")
		panic("")
	}
	cleanUp()
	os.Exit(0)
}

func buildBinary() string {
	_, err := runTestCommand("/bin/sh", filepath.Join(".", "features", "build-cli.sh"))
	if err != nil {
		panic(err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(".", "features", "test-platform.exe")
	case "ios":
		return filepath.Join(".", "features", "test-platform-macos")
	case "linux":
		return filepath.Join(".", "features", "test-platform-linux")
	default:
		panic(errors.New("Operating system not supported"))
	}
}

func copyFiles() {
	_, err := runTestCommand("/bin/sh", filepath.Join(".", "features", "copy-files.sh"))
	if err != nil {
		panic(err)
	}
}

func runTestCommand(commandName string, commandSlice ...string) (string, error) {
	cmd := exec.Command(commandName, commandSlice...)
	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "Error creating stdOutPipe for Cmd.")
	}
	stdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.Wrap(err, "Error creating stdErrPipe for Cmd.")
	}

	var output []rune
	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			s := stdOutScanner.Text()
			if s != "" {
				for _, ss := range s {
					output = append(output, ss)
				}
				output = append(output, '\n')
			}
		}
	}()

	var errStr []rune
	stdErrScanner := bufio.NewScanner(stdErrReader)
	go func() {
		for stdErrScanner.Scan() {
			s := stdErrScanner.Text()
			if s != "" {
				for _, ss := range s {
					errStr = append(errStr, ss)
				}
				errStr = append(errStr, '\n')
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, string(errStr))
	}

	if string(errStr) != "" {
		return "", errors.Wrap(errors.New(string(errStr)), "")
	}

	return string(output), nil
}

func deleteContentAtFilePath(filePath []string, content []string) {
	for _, c := range content {
		err := os.RemoveAll(filepath.Join(filepath.Join(filePath...), c))
		if err != nil {
			panic(err)
		}
	}
}

func cleanUp() {
	deleteContentAtFilePath([]string{".", "features"}, []string{"test-platform.exe", "test-platform-linux", "test-platform-macos"})

	_, err := runTestCommand("docker", "rm", "instant-openhie")
	if err == nil {
		fmt.Println("[ERROR]: ", err)
	}

	_, err = runTestCommand("docker", "volume", "rm", "instant")
	if err != nil && !strings.Contains(err.Error(), "No such volume: instant") {
		fmt.Println("[ERROR]: ", err)
	}

	err = os.Remove(".env.test")
	if err != nil {
		fmt.Println("[ERROR]:", err)
	}

	err = os.Remove("config.yml")
	if err != nil {
		fmt.Println("[ERROR]:", err)
	}
}

func hasCustomPackage(command string) ([]string, bool) {
	if strings.Contains(command, "-c=") {
		split := strings.SplitAfter(command, "-c=")

		var customPackageNames []string
		for i := 1; i < len(split); i++ {
			subSplit := strings.Split(split[1], " ")
			customPackageNames = append(customPackageNames, strings.TrimSuffix(path.Base(path.Clean(subSplit[0])), path.Ext(subSplit[0])))
		}

		return customPackageNames, true
	} else if strings.Contains(command, "--custom-path=") {
		split := strings.SplitAfter(command, "--custom-path=")

		var customPackageNames []string
		for i := 1; i < len(split); i++ {
			subSplit := strings.Split(split[1], " ")
			customPackageNames = append(customPackageNames, strings.TrimSuffix(path.Base(path.Clean(subSplit[0])), path.Ext(subSplit[0])))
		}

		return customPackageNames, true
	}

	return nil, false
}

// This function can run infinitely, the caller must ensure that the goroutine running
// this process is terminated
func monitorDirFor(directory, basename string) {
	pathName := filepath.Join(directory, basename)
	for {
		_, err := os.Stat(pathName)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			os.Exit(1)
		}

		customPackages[basename] = true
		break
	}
}
