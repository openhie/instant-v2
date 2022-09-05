package install

import (
	v1 "github.com/openhie/package-starter-kit/cli"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func InitNoneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "none",
		Short: "Install with no auth",
		Run: func(cmd *cobra.Command, args []string) {
			params := &v1.Params{}
			params.TypeAuth = "None"

			urlEntry, err := cmd.Flags().GetString("url-entry")
			util.LogError(err)
			fhirServer, err := cmd.Flags().GetString("fhir-server")
			util.LogError(err)

			err = v1.LoadIGpackage(urlEntry, fhirServer, params)
			util.LogError(err)
		},
	}

	flags := cmd.Flags()

	flags.StringP("url-entry", "url", "", "The url entry for the fhir IG")
	flags.StringP("fhir-server", "fhir", "", "The fhir server for the IG")

	return cmd
}
