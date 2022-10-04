package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DeclareConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config management commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("To be implemented")
		},
	}
	return cmd
}
