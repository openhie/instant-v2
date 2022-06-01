package pkg

import (
	"github.com/spf13/cobra"
)

func InitPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	cmd.AddCommand(
		InitInitCommand(),
		InitUpCommand(),
		InitDownCommand(),
		InitRemoveCommand(),
	)

	return cmd
}
