package cmd

import (
	"context"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"

	"cli/cmd/commands"
	"cli/core/state"
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
	rootCmd.PersistentFlags().StringVar(&state.ConfigFile, "config", "", "config file (default is $WORKING_DIR/config.yaml)")
	// Note: No shorthand for env-file, saving -e for individual env var declarations
	rootCmd.PersistentFlags().StringSliceVar(&state.EnvFiles, "env-file", nil, "env file")

	commands.AddCommands(rootCmd)
}
