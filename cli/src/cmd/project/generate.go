package project

import (
	"github.com/spf13/cobra"
)

// TODO(MarkL): Write tests for this once this functionality is introduced
func ProjectGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return cmd
}
