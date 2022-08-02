package config

import (
	"fmt"

	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/spf13/cobra"
)

func DeclareConfigCommand(global *types.Global) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config management commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("To be implemented")
		},
	}
	return cmd
}
