package main

import (
	"embed"

	"ohiecli/cmd"
	"ohiecli/core"
)

//go:embed template/*
var templateFs embed.FS

func main() {
	core.TemplateFs = templateFs
	cmd.Execute()
}
