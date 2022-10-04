package main

import (
	"embed"

	"cli/cmd"
	"cli/core"
)

//go:embed template/*
var templateFs embed.FS

func main() {
	core.TemplateFs = templateFs
	cmd.Execute()
}
