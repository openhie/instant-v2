package completion

import (
	"context"

	"cli/core/parse"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func FlagCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := parse.GetConfigFromParams(cmd)
		if err != nil {
			log.Error(context.Background(), err)
		}

		return config.Packages, cobra.ShellCompDirectiveDefault
	})
	cmd.RegisterFlagCompletionFunc("profile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := parse.GetConfigFromParams(cmd)
		if err != nil {
			log.Error(context.Background(), err)
		}

		var profileNames []string
		for _, p := range config.Profiles {
			profileNames = append(profileNames, p.Name)
		}

		return profileNames, cobra.ShellCompDirectiveDefault
	})
	cmd.RegisterFlagCompletionFunc("custom-path", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		config, err := parse.GetConfigFromParams(cmd)
		if err != nil {
			log.Error(context.Background(), err)
		}

		var customPackages []string
		for _, c := range config.CustomPackages {
			customPackages = append(customPackages, c.Id)
		}

		return customPackages, cobra.ShellCompDirectiveDefault
	})
}

func GenCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate the autocompletion script for the specified shell",
	}

	cmd.AddCommand(
		genBashCompletionCommand(),
		genZshCompletionCommand(),
	)

	return cmd
}
