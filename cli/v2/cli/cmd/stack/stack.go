package stack

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/spf13/cobra"
)

func DeclareStackCommand(global *types.Global) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Stack level commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("To be implemented")
		},
	}
	return cmd
}
