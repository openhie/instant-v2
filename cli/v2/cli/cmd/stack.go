package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stackCmd represents the stack command
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Stack level commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}

func init() {
	rootCmd.AddCommand(stackCmd)
}
