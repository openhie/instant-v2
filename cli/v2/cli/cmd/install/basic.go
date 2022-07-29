package install

import (
	v1 "github.com/openhie/package-starter-kit/cli"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitBasicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basic",
		Short: "Install with basic auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.Params{}
			params.TypeAuth = "Basic"

			basicUser, err := cmd.Flags().GetString("basic-user")
			util.LogError(err)
			basicPassword, err := cmd.Flags().GetString("basic-password")
			util.LogError(err)
			params.BasicUser = basicUser
			params.BasicPass = basicPassword

			err = v1.loadIGpackage(startupCommands[1], startupCommands[2], params)
			util.LogError(err)
		},
	}

	flags := cmd.Flags()

	flags.StringP("basic-user", "bu", "", "The basic user credential for auth")
	flags.StringP("basic-password", "bp", "", "The basic password credential for auth")

	return cmd
}
