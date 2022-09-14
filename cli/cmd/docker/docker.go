package docker

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"ohie_cli/config"
	"ohie_cli/utils"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var (
	OsCreate           = os.Create
	IoCopy             = io.Copy
	ZipOpenReader      = zip.OpenReader
	OsMkdirAll         = os.MkdirAll
	FilepathJoin       = filepath.Join
	OsOpenFile         = os.OpenFile
	OsRemove           = os.Remove
	execCommand        = exec.Command
	RunDeployCommand   = runDeployCommand
	RunCommand         = runCommand
	MountCustomPackage = mountCustomPackage
)

type CommandsOptions struct {
	environmentVariables []string
	deployCommand        string
	otherFlags           []string
	packages             []string
	customPackagePaths   []string
	imageVersion         string
	targetLauncher       string
}

var logsToSuppress = []string{
	"Already exists",
	"Pulling fs layer",
	"Verifying Checksum",
	"Download complete",
	"Waiting",
	"Pull complete",
}

func DebugDocker() error {
	fmt.Printf("...checking your Docker setup")

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Can't get current working directory")
	}

	fmt.Println(cwd)

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to get Docker context. Please ensure that Docker is downloaded and running")
	} else {
		// Docker default is 2GB, which may need to be revisited if Instant grows.
		str1 := "bytes memory is allocated\n"
		str2 := strconv.FormatInt(info.MemTotal, 10)
		result := str2 + str1
		fmt.Println(result)
		fmt.Println("Docker setup looks good")
	}

	return nil
}

func getPackagePaths(inputArr []string, flags []string) (packagePaths []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				packagePath := strings.Replace(i, flag, "", 1)
				packagePath = strings.Trim(packagePath, "\"")
				packagePaths = append(packagePaths, packagePath)
			}
		}
	}
	return
}

func getEnvironmentVariables(inputArr []string, flags []string) (environmentVariables []string) {
	for _, i := range inputArr {
		for _, flag := range flags {
			if strings.Contains(i, flag) {
				environmentVariables = append(environmentVariables, strings.SplitN(i, "=", 2)...)
			}
		}
	}
	return
}

func extractCommands(startupCommands []string) CommandsOptions {
	imageVersion := "latest"

	if strings.Contains(config.Cfg.Image, ":") {
		imageVersion = strings.Split(config.Cfg.Image, ":")[1]
	}

	commandOptions := CommandsOptions{
		imageVersion: imageVersion,
	}

	for _, option := range startupCommands {
		switch {
		case utils.SliceContains([]string{"init", "up", "down", "destroy"}, option):
			commandOptions.deployCommand = option
		case strings.HasPrefix(option, "-c=") || strings.HasPrefix(option, "--custom-package="):
			commandOptions.customPackagePaths = append(commandOptions.customPackagePaths, option)
		case strings.HasPrefix(option, "-e=") || strings.HasPrefix(option, "--env-file="):
			commandOptions.environmentVariables = append(commandOptions.environmentVariables, option)
		case strings.HasPrefix(option, "--image-version="):
			commandOptions.imageVersion = strings.Split(option, "--image-version=")[1]
		case strings.HasPrefix(option, "-t="):
			commandOptions.targetLauncher = strings.Split(option, "-t=")[1]
		case strings.HasPrefix(option, "-") || strings.HasPrefix(option, "--"):
			commandOptions.otherFlags = append(commandOptions.otherFlags, option)
		default:
			commandOptions.packages = append(commandOptions.packages, option)
		}
	}

	if len(commandOptions.customPackagePaths) > 0 {
		commandOptions.customPackagePaths = getPackagePaths(commandOptions.customPackagePaths, []string{"-c=", "--custom-package="})
	}

	if len(commandOptions.environmentVariables) > 0 {
		commandOptions.environmentVariables = getEnvironmentVariables(commandOptions.environmentVariables, []string{"-e=", "--env-file="})
	}

	if commandOptions.targetLauncher == "" {
		commandOptions.targetLauncher = config.CustomOptions.TargetLauncher
	}

	return commandOptions
}

func runDeployCommand(startupCommands []string) error {
	fmt.Println("Note: Initial setup takes 1-5 minutes.\nWait for the DONE message.\n--------------------------")

	commandOptions := extractCommands(startupCommands)

	if len(commandOptions.packages) == 0 {
		for _, p := range config.Cfg.Packages {
			commandOptions.packages = append(commandOptions.packages, p.ID)
		}
	}

	fmt.Println("Action:", commandOptions.deployCommand)
	fmt.Println("Package IDs:", commandOptions.packages)
	fmt.Println("Custom package paths:", commandOptions.customPackagePaths)
	fmt.Println("Environment Variables:", commandOptions.environmentVariables)
	fmt.Println("Other Flags:", commandOptions.otherFlags)
	fmt.Println("Image Version:", commandOptions.imageVersion)
	fmt.Println("Target Launcher:", commandOptions.targetLauncher)

	image := ""
	if strings.Contains(config.Cfg.Image, ":") {
		image = strings.Split(config.Cfg.Image, ":")[0] + ":" + commandOptions.imageVersion
	} else {
		image = config.Cfg.Image + ":" + commandOptions.imageVersion
	}

	fmt.Println("Creating fresh instant container with volumes...")
	commandSlice := []string{
		"create",
		"--rm",
		"--mount=type=volume,src=instant,dst=/instant",
		"--name", "instant-openhie",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"--network", "host",
	}

	if config.Cfg.LogPath != "" {
		commandSlice = append(commandSlice, fmt.Sprintf("--mount=type=bind,src=%s,dst=/tmp/logs", config.Cfg.LogPath))
	}

	commandSlice = append(commandSlice, commandOptions.environmentVariables...)
	commandSlice = append(commandSlice, []string{image, commandOptions.deployCommand}...)
	commandSlice = append(commandSlice, commandOptions.otherFlags...)
	commandSlice = append(commandSlice, []string{"-t", commandOptions.targetLauncher}...)
	commandSlice = append(commandSlice, commandOptions.packages...)

	_, err := RunCommand("docker", nil, commandSlice...)
	if err != nil {
		return err
	}
	defer removeInstantVolume()

	fmt.Println("Adding 3rd party packages to instant volume:")

	for _, c := range commandOptions.customPackagePaths {
		fmt.Print("- " + c)
		err = MountCustomPackage(c)
		if err != nil {
			return err
		}
	}

	fmt.Println("\nRunning orchestration container")
	commandSlice = []string{"start", "-a", "instant-openhie"}
	_, err = RunCommand("docker", nil, commandSlice...)
	if err != nil {
		color.Red("\nError: Failed while running orchestration container, check output from Orchestration container above. (Underlying error: " + err.Error() + ")")
		// ignore error and return user to prompt
		return nil
	}

	return nil
}

func removeInstantVolume() {
	fmt.Println("\n\nRemoving instant volume...")
	_, err := RunCommand("docker", []string{"Error: No such volume: instant"}, []string{"volume", "rm", "instant"}...)
	if err != nil {
		fmt.Println(errors.Wrap(err, "[Error] Failed to remove instant volume."))
	}
}

var runCommand = func(commandName string, suppressErrors []string, commandSlice ...string) (pathToPackage string, err error) {
	cmd := execCommand(commandName, commandSlice...)
	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return pathToPackage, errors.Wrap(err, "Error creating stdOutPipe for Cmd.")
	}
	stdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return pathToPackage, errors.Wrap(err, "Error creating stdErrPipe for Cmd.")
	}

	messages := make(chan string)
	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			messages <- fmt.Sprintf("\t > %s", stdOutScanner.Text())
		}
	}()

	var stderr string
	stdErrScanner := bufio.NewScanner(stdErrReader)
	go func() {
		for stdErrScanner.Scan() {
			if stdErrScanner.Text() != "" {
				stderr = stdErrScanner.Text()
				if utils.SliceContainsSubstr(logsToSuppress, stderr) {
				} else if utils.SliceContainsSubstr([]string{
					"Unable to find image",
					"Pulling from",
					"Downloaded newer image",
					"Digest",
				}, stderr) {
					messages <- fmt.Sprintf("\t > %s", stderr)
				} else {
					messages <- fmt.Sprintf("\t > [ERROR] %s", stderr)
				}
			}
		}
	}()

	go func() {
		for message := range messages {
			if strings.Contains(message, "ERROR") {
				color.Red("%s", message)
			} else {
				fmt.Println(message)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		if suppressErrors != nil && utils.SliceContains(suppressErrors, stderr) {
		} else {
			return pathToPackage, errors.Wrap(err, "Error starting Cmd. "+stderr)
		}
	}

	err = cmd.Wait()
	if err != nil {
		if suppressErrors != nil && utils.SliceContains(suppressErrors, stderr) {
		} else {
			return pathToPackage, errors.Wrap(err, "Error waiting for Cmd. "+stderr)
		}
	}

	if commandName == "git" {
		if len(commandSlice) < 2 {
			return pathToPackage, errors.New("Not enough arguments for git command")
		}
		pathToPackage = commandSlice[1]
		// Get name of repo
		urlSplit := strings.Split(pathToPackage, ".")
		urlPathSplit := strings.Split(urlSplit[len(urlSplit)-2], "/")
		repoName := urlPathSplit[len(urlPathSplit)-1]

		pathToPackage = filepath.Join(".", repoName)
	}

	return pathToPackage, nil
}

func mountCustomPackage(pathToPackage string) error {
	gitRegex := regexp.MustCompile(`\.git`)
	httpRegex := regexp.MustCompile("http")
	zipRegex := regexp.MustCompile(`\.zip`)
	tarRegex := regexp.MustCompile(`\.tar`)

	var err error
	if gitRegex.MatchString(pathToPackage) {
		pathToPackage, err = runCommand("git", nil, []string{"clone", pathToPackage}...)
		if err != nil {
			return err
		}
	} else if httpRegex.MatchString(pathToPackage) {
		resp, err := http.Get(pathToPackage)
		if err != nil {
			return errors.Wrap(err, "Error in downloading custom package")
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errors.Wrapf(err, "Error in downloading custom package - HTTP status code: %v", strconv.Itoa(resp.StatusCode))
		}

		if zipRegex.MatchString(pathToPackage) {
			pathToPackage, err = unzipPackage(resp.Body)
			if err != nil {
				return err
			}
		} else if tarRegex.MatchString(pathToPackage) {
			pathToPackage, err = untarPackage(resp.Body)
			if err != nil {
				return err
			}
		}
	}

	commandSlice := []string{"cp", pathToPackage, "instant-openhie:instant/"}
	_, err = runCommand("docker", nil, commandSlice...)
	return err
}

func createZipFile(file string, content io.Reader) error {
	output, err := OsCreate(file)
	if err != nil {
		return errors.Wrap(err, "Error in creating zip file:")
	}
	defer output.Close()

	bytesWritten, err := IoCopy(output, content)
	if err != nil {
		return errors.Wrap(err, "Error in copying zip file content:")
	}
	if bytesWritten < 1 {
		return errors.New("File created but no content written.")
	}

	return nil
}

var unzipPackage = func(zipContent io.ReadCloser) (pathToPackage string, err error) {
	tempZipFile := "temp.zip"
	err = createZipFile(tempZipFile, zipContent)
	if err != nil {
		return "", err
	}

	// Unzip file
	archive, err := ZipOpenReader(tempZipFile)
	if err != nil {
		return "", errors.Wrap(err, "Error in unzipping file:")
	}
	defer archive.Close()

	packageName := ""
	for _, file := range archive.File {
		filePath := FilepathJoin(".", file.Name)

		if file.FileInfo().IsDir() {
			if packageName == "" {
				packageName = file.Name
			}
			err = OsMkdirAll(filePath, os.ModePerm)
			if err != nil {
				return "", err
			}
			continue
		}

		content, err := file.Open()
		if err != nil {
			return "", errors.Wrap(err, "Error in unzipping file:")
		}
		defer content.Close()

		dest, err := OsOpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", errors.Wrap(err, "Error in unzipping file:")
		}
		defer dest.Close()

		written, err := IoCopy(dest, content)
		if err != nil {
			return "", errors.Wrap(err, "Error in copying unzipping file:")
		}
		if written < 1 {
			return "", errors.New("No content copied")
		}
	}

	// Remove temp zip file
	tempFilePath := FilepathJoin(".", tempZipFile)
	archive.Close()
	err = OsRemove(tempFilePath)
	if err != nil {
		return "", errors.Wrap(err, "Error in deleting temp.zip file:")
	}

	return FilepathJoin(".", packageName), nil
}

var untarPackage = func(tarContent io.ReadCloser) (pathToPackage string, err error) {
	packageName := ""
	gzipReader, err := gzip.NewReader(tarContent)
	if err != nil {
		return "", errors.Wrap(err, "Error in extracting tar file")
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if file == nil {
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "Error in extracting tar file")
		}

		filePath := filepath.Join(".", file.Name)
		if file.Typeflag == tar.TypeDir {
			if packageName == "" {
				packageName = filePath
			}
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return "", errors.Wrap(err, "Error in untaring file")
		}

		_, err = io.Copy(dest, tarReader)
		if err != nil {
			return "", errors.Wrap(err, "Error in extracting tar file")
		}
	}
	pathToPackage = filepath.Join(".", packageName)

	return pathToPackage, nil
}

// TODO: Check references
func StopContainer() {
	commandSlice := []string{"stop", "instant-openhie"}
	suppressErrors := []string{"Error response from daemon: No such container: instant-openhie"}
	_, err := RunCommand("docker", suppressErrors, commandSlice...)
	if err != nil {
		log.Fatalf("runCommand() failed: %v", err)
	}
}

// Gracefully shut down the instant container and then kill the go cli with the panic error or message passed.
func GracefulPanic(err error, message string) {
	StopContainer()
	if message != "" {
		panic(message)
	}
	panic(err)
}
