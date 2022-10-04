package install

import (
	"ohiecli/old/config"

	"github.com/spf13/cobra"
)

func InitCustomCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom",
		Short: "Install with custom auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := config.Params{}
			params.TypeAuth = "Custom"
		},
	}

	return cmd
}
