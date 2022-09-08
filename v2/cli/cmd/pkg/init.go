package pkg

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/core"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
)

func PackageInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a package with relevant configs, volumes and setup",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := getConfigFromParams(cmd)
			util.PanicError(err)
			packageSpec, err := getPackageSpecFromParams(cmd)
			util.PanicError(err)
			packageSpec, err = loadInProfileParams(cmd, *config, *packageSpec)
			util.PanicError(err)

			sshKey, err := cmd.Flags().GetString("ssh-key")
			util.PanicError(err)
			packageSpec.SSHKeyFile = sshKey

			sshPassword, err := cmd.Flags().GetString("ssh-password")
			util.PanicError(err)
			packageSpec.SSHPasswordFile = sshPassword

			err = core.LaunchPackage(*packageSpec, *config)
			util.PanicError(err)
		},
	}

	setPackageActionFlags(cmd)
	flags := cmd.Flags()
	flags.String("ssh-key", "", "The path to the ssh key required for cloning a custom package")
	flags.String("ssh-password", "", "The password (or path to the file containing the password) required for authenticating the ssh-key when cloning a custom package")

	return cmd
}
