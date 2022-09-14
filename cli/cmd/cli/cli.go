package cli

import (
	"errors"
	"fmt"
	"ohie_cli/config"
	"ohie_cli/docker"
	"ohie_cli/ig"
	"ohie_cli/utils"
	"os"
)

func CLI() error {
	startupCommands := os.Args[1:]

	var err error
	switch startupCommands[0] {
	case "help":
		fmt.Println(utils.GetHelpText(false, ""))
	case "version":
	case "install":
		params := &config.Params{}
		switch startupCommands[3] {
		case "none", "None":
			params.TypeAuth = "None"
			err = ig.LoadIGpackage(startupCommands[1], startupCommands[2], params)
		case "basic", "Basic":
			params.TypeAuth = "Basic"
			params.BasicUser = startupCommands[4]
			params.BasicPass = startupCommands[5]
			err = ig.LoadIGpackage(startupCommands[1], startupCommands[2], params)
		case "token", "Token":
			params.TypeAuth = "Token"
			params.Token = startupCommands[4]
			err = ig.LoadIGpackage(startupCommands[1], startupCommands[2], params)
		case "custom", "Custom":
			params.TypeAuth = "Custom"
			params.Token = startupCommands[4]
			err = ig.LoadIGpackage(startupCommands[1], startupCommands[2], params)
		}
	default:
		if len(startupCommands) < 2 {
			fmt.Println("The deploy command is not recognized: ", startupCommands)
			return errors.New("incorrect arguments list passed to CLI. Requires at least 2 arguments when in non-interactive mode")
		}

		err = docker.RunDeployCommand(startupCommands)
		if err != nil {
			return err
		}
	}

	return err
}
