package install

import (
	v1 "github.com/openhie/package-starter-kit/cli"
	"github.com/spf13/cobra"
)

func InitTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Install with token auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.Params{}
			params.TypeAuth = "Token"
		},
	}

	return cmd
}
