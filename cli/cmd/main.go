package main

import (
	"embed"
	"io/ioutil"
	"log"
	"ohie_cli/cli"
	"ohie_cli/config"
	"ohie_cli/docker"
	"ohie_cli/prompts"
	"os"
	"strings"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v3"
)

//go:embed banner.txt
//go:embed version
var f embed.FS

func LoadConfig() {
	yamlConfig, loadErr := ioutil.ReadFile("config.yml")
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	err := yaml.Unmarshal(yamlConfig, &config.Cfg)
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

func main() {
	LoadConfig()
	showBanner()

	//Need to set the default here as we declare the struct before the config is loaded in.
	config.CustomOptions.TargetLauncher = config.Cfg.DefaultTargetLauncher

	if strings.Contains(config.Cfg.Image, ":") {
		config.CustomOptions.ImageVersion = strings.Split(config.Cfg.Image, ":")[1]
	}

	version, err := f.ReadFile("version")
	if err != nil {
		log.Println(err)
	}

	color.Cyan("Go Cli Version: " + string(version))
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	if len(os.Args) > 1 {
		err = cli.CLI()
		if err != nil {
			docker.GracefulPanic(err, "")
		}
	} else {
		err = prompts.SelectSetup()
		if err != nil {
			docker.GracefulPanic(err, "")
		}
	}
}
