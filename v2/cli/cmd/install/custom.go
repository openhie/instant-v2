package install

import (
	v1 "github.com/openhie/package-starter-kit/cli"
	"github.com/spf13/cobra"
)

func InitCustomCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom",
		Short: "Install with custom auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.Params{}
			params.TypeAuth = "Custom"
		},
	}

	return cmd
}
