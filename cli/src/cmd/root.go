package cmd

import (
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"

	"cli/cmd/commands"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A cli to assist with package deployment and management",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error(context.Background(), err)
		panic(err)
	}
}

func init() {
	// TODO: read the docs for cobra.OnInitialize() and decide if it's needed
	cobra.OnInitialize()

	commands.AddCommands(rootCmd)
}
