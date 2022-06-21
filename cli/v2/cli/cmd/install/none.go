package install

import (
	v1 "github.com/openhie/package-starter-kit/cli"
	"github.com/spf13/cobra"
)

func InitNoneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "none",
		Short: "Install with no auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.Params{}
			params.TypeAuth = "None"
		},
	}

	return cmd
}
