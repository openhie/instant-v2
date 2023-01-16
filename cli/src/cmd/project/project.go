package project

import (
	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
)

func checkInvalidFlags(cmd *cobra.Command) error {
	if cmd.Flag("name").Changed {
		return errors.Wrap(errors.New("flag accessed but not defined: name"), "")
	}
	if cmd.Flag("profile").Changed {
		return errors.Wrap(errors.New("flag accessed but not defined: profile"), "")
	}

	return nil
}

func DeclareProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project level commands",
	}

	cmd.AddCommand(
		projectInitCommand(),
		projectDownCommand(),
		projectUpCommand(),
		projectDestroyCommand(),
		projectGenerateCommand(),
	)

	return cmd
}
