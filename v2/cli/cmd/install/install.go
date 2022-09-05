package install

import (
	"github.com/spf13/cobra"
)

func InitInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Package level commands",
	}

	cmd.AddCommand(
		InitNoneCommand(),
		InitBasicCommand(),
		InitTokenCommand(),
		InitCustomCommand(),
	)

	return cmd
}
