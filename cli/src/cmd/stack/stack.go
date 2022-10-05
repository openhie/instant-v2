package stack

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DeclareStackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Stack level commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("To be implemented")
		},
	}
	
	return cmd
}
