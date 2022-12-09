package pkg

import (
	"github.com/spf13/cobra"
)

func DeclarePackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Package level commands",
	}

	cmd.AddCommand(
		packageInitCommand(),
		packageUpCommand(),
		packageDownCommand(),
		packageRemoveCommand(),
		packageGenerateCommand(),
	)

	return cmd
}
