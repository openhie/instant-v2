package main

import (
	"embed"
	"log"
	"os"
	"github.com/openhie/package-starter-kit/cli/pkg"
	"github.com/fatih/color"
)

//go:embed banner.txt
//go:embed version
var f embed.FS

func main() {
	pkg.loadConfig()
	pkg.showBanner()

	//Need to set the default here as we declare the struct before the config is loaded in.
	customOptions.targetLauncher = pkg.cfg.DefaultTargetLauncher

	version, err := f.ReadFile("version")
	if err != nil {
		log.Println(err)
	}

	color.Cyan("Version: " + string(version))
	color.Blue("Remember to stop applications or they will continue to run and have an adverse impact on performance.")

	if len(os.Args) > 1 {
		err = CLI()
		if err != nil {
			gracefulPanic(err, "")
		}
	} else {
		err = selectSetup()
		if err != nil {
			gracefulPanic(err, "")
		}
	}
}
