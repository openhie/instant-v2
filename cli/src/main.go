package main

import (
	"context"
	"embed"
	"os"

	"cli/cmd"
	"cli/core"

	"github.com/docker/docker/client"
	"github.com/luno/jettison/log"
)

//go:embed template/*
var templateFs embed.FS

func main() {
	defer handleExit()

	core.TemplateFs = templateFs
	cmd.Execute()
}

func handleExit() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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
