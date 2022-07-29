package main

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/openhie/package-starter-kit/cli/pkg"
	"gopkg.in/yaml.v2"
)

//go:embed banner.txt
//go:embed version
var f embed.FS

var cfg Config

type Package struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type Config struct {
	Image                        string    `yaml:"image"`
	DefaultTargetLauncher        string    `yaml:"defaultTargetLauncher"`
	Packages                     []Package `yaml:"packages"`
	DisableKubernetes            bool      `yaml:"disableKubernetes"`
	DisableIG                    bool      `yaml:"disableIG"`
	DisableCustomTargetSelection bool      `yaml:"disableCustomTargetSelection"`
	LogPath                      string    `yaml:"logPath"`
}

type customOption struct {
	startupAction              string
	startupPackages            []string
	envVarFileLocation         string
	envVars                    []string
	customPackageFileLocations []string
	onlyFlag                   bool
	imageVersion               string
	targetLauncher             string
	devMode                    bool
}

var customOptions = customOption{
	startupAction:      "init",
	envVarFileLocation: "",
	onlyFlag:           false,
	imageVersion:       "latest",
	targetLauncher:     "docker",
	devMode:            false,
}

func stopContainer() {
	commandSlice := []string{"stop", "instant-openhie"}
	suppressErrors := []string{"Error response from daemon: No such container: instant-openhie"}
	_, err := pkg.RunCommand("docker", suppressErrors, commandSlice...)
	if err != nil {
		log.Fatalf("runCommand() failed: %v", err)
	}
}

//Gracefully shut down the instant container and then kill the go cli with the panic error or message passed.
func gracefulPanic(err error, message string) {
	stopContainer()
	if message != "" {
		panic(message)
	}
	panic(err)
}

func loadConfig() {
	yamlConfig, loadErr := ioutil.ReadFile("config.yml")
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	err := yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func showBanner() {
	// Check for custom banner, otherwise use embedded
	banner, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		banner, err = f.ReadFile("banner.txt")
		if err != nil {
			log.Println(err)
		}
	}

	color.Green(string(banner))
}

func getHelpText(interactive bool, options string) string {
	if interactive {
		switch options {
		case "Deploy Commands":
			return `Commands:
				init/up/destroy/down	the deploy command you want to run (brief description below)
					deploy commands:
						init:			for initializing a service
						up:				for starting up a service that has been shut down or updating a service
						destroy:	for destroying a service
						down:			for bringing down a running service
			`
		case "Custom Options":
			return `Commands:
				Choose deploy action - for choosing the deploy action

				Choose target launcher - for choosing the deploy target. Can be a docker swarm, kubernetes or docker (Project specific)

				Specify deploy packages - for choosing the packages you want to use (core and custom packages)

				Specify environment variable file location - for specifying the file path to an environment variables file

				Specify environment variables - for specifying environment variables

				Specify custom package locations - for specifying the location or url to the custom packages you want to operate on

				Toggle only flag - for specifying the only flag, which specifies that actions are to be taken on a single package and not on its dependencies

				Specify Image Version - for specifying the version of the instant or platform image to use. Default is latest

				Toggle dev mode - for enabling the development mode in which the service ports are exposed

				Execute with current options - this executes the options that have been specified

				View current options set - for viewing the options that have been specified

				Reset to default options - for resetting to default options
			`
		default:
			return `Commands:
				Use Docker on your PC - this is for deploying packages to either docker or docker swarm

				Use a kubernetes Cluster - this is for deploying packages to a kubernetes cluster

				Install FHIR package - this is for installing FHIR IGs hosted remotely
			`
		}
	} else {
		return `Commands: 
		help 		this menu

		init/up/destroy/down	the deploy command you want to run (brief description below)
					deploy commands:
						init:	 for initializing a service
						up:	 for starting up a service that has been shut down or updating a service
						destroy: for destroying a service
						down:	 for bringing down a running service
					custom flags:
						--only, -o:							used to specify a single service for services that have dependencies. For cases where one wants to shut down or destroy a service without affecting its dependencies
						-t:											specifies the target to deploy to. Options are docker, swarm (docker swarm) and k8s (kubernetes) - project dependant.
						--custom-package, -c:		specifies path or url to a custom package. Git ssh urls are supported
						--dev:									specifies the development mode in which all service ports are exposed
						-e:											for specifying an environment variable
						--env-file: 						for specifying the path to an environment variables file
						--image-version:			the version of the project used for the deploy. Defaults to 'latest'
						-*, --*:								unrecognised flags are passed through uninterpreted
					usage:
						<deploy command> <custom flags> <package ids>
					examples:
						{your_binary_file} init -t=swarm --dev -e="NODE_ENV=prod" --env-file="../env.dev" -c="../customPackage1" -c="<git@github.com/customPackage2>"  interoperability-layer-openhim customPackage1_id customPackage2_id
						{your_binary_file} down -t=docker --only elastic_analytics

		install		install fhir npm package on fhir server
					usage: install <ig_url> <fhir_server> <authtype> <user/token> <pass>

					examples:
					install https://intrahealth.github.io/simple-hiv-ig/ http://hapi.fhir.org/baseR4 none
					install <ig_url> <fhir_server> basic smith stuff
					install <ig_url> <fhir_server> token "123"
					install <ig_url> <fhir_server> custom test
		`
	}
}

func main() {
	loadConfig()
	showBanner()

	//Need to set the default here as we declare the struct before the config is loaded in.
	customOptions.targetLauncher = cfg.DefaultTargetLauncher

	if strings.Contains(cfg.Image, ":") {
		customOptions.imageVersion = strings.Split(cfg.Image, ":")[1]
	}

	version, err := f.ReadFile("version")
	if err != nil {
		log.Println(err)
	}

	color.Cyan("Go Cli Version: " + string(version))
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	if len(os.Args) > 1 {
		err = pkg.CLI()
		if err != nil {
			gracefulPanic(err, "")
		}
	} else {
		err = pkg.SelectSetup()
		if err != nil {
			gracefulPanic(err, "")
		}
	}
}
