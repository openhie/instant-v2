package main

import (
	"embed"

	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd"
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
)

//go:embed template/*
var templateFs embed.FS

func main() {
	core.TemplateFs = templateFs
	cmd.Execute()
}
