package main

import (
	"context"
	"embed"
	"os"

	"cli/cmd"
	"cli/core"
	"cli/util"

	"github.com/luno/jettison/log"
)

//go:embed template/*
var templateFs embed.FS

// TODO: starting using context correctly
func main() {
	defer handleExit()

	core.TemplateFs = templateFs
	cmd.Execute()
}

// TODO: should the cli panic if launching a package not contained in the config.yml?
func handleExit() {
	ctx := context.Background()

	cli, err := util.NewDockerClient()
	if err != nil {
		log.Error(ctx, err)
	}

	core.RemoveStaleInstantContainer(cli, ctx)
	core.RemoveStaleInstantVolume(cli, ctx)

	if recover() != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
