package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package level commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}
