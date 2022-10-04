package cmd

import (
	"os"

	"cli/cmd/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var configViper viper.Viper

var envFiles []string
var envVarViper viper.Viper

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
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $WORKING_DIR/config.yaml)")
	rootCmd.PersistentFlags().StringSliceVarP(&envFiles, "env-file", "e", nil, "env file (default is $WORKING_DIR/.env)")

	commands.AddCommands(rootCmd)
}
