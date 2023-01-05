package main

import (
	"context"
	"os"

	"cli/cmd"
	"cli/util/docker"

	"github.com/luno/jettison/log"
)

func main() {
	defer handleExit()

	cmd.Execute()
}

func handleExit() {
	ctx := context.Background()

	cli, err := docker.NewDockerClient()
	if err != nil {
		log.Error(ctx, err)
	}

	docker.RemoveStaleInstantContainer(cli, ctx)
	docker.RemoveStaleInstantVolume(cli, ctx)

	if recover() != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
