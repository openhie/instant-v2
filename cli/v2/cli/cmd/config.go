package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config management commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
