package pkg

import (
	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

var ErrNoPackages = errors.New("no packages selected in any of command-line/profiles, use the 'project' command for project level functions")

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
