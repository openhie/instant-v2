package install

import (
	"ohiecli/util"

	"github.com/spf13/cobra"
)

func InitBasicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "basic",
		Short: "Install with basic auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := params{}
			params.TypeAuth = "Basic"

			basicUser, err := cmd.Flags().GetString("basic-user")
			util.LogError(err)
			basicPassword, err := cmd.Flags().GetString("basic-password")
			util.LogError(err)
			params.BasicUser = basicUser
			params.BasicPass = basicPassword

			urlEntry, err := cmd.Flags().GetString("url-entry")
			util.LogError(err)
			fhirServer, err := cmd.Flags().GetString("fhir-server")
			util.LogError(err)

			err = loadIGpackage(urlEntry, fhirServer, &params)
			util.LogError(err)
		},
	}

	flags := cmd.Flags()

	flags.StringP("basic-user", "bu", "", "The basic user credential for auth")
	flags.StringP("basic-password", "bp", "", "The basic password credential for auth")

	return cmd
}
