package pkg

import (
	"github.com/spf13/cobra"
)

// TODO(MarkL): Write tests for this once this functionality is introduced
func packageGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate a new package",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	return cmd
}
