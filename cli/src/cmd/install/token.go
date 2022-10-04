package install

import (
	"github.com/spf13/cobra"
)

func InitTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Install with token auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := params{}
			params.TypeAuth = "Token"
		},
	}

	return cmd
}
