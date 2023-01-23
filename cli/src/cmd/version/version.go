package version

import (
	"context"
	"embed"
	"fmt"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

//go:embed version
var versionFile embed.FS

func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			version, err := versionFile.ReadFile("version")
			if err != nil {
				log.Error(context.Background(), err)
			}

			fmt.Println("instant-CLI version:", string(version))
		},
	}

	return cmd
}
