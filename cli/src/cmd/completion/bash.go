package completion

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	execShell "cli/util/exec"
	fileUtil "cli/util/file"

	"github.com/luno/jettison/log"
	"github.com/spf13/cobra"
)

func genBashCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bash",
		Short: "Generate the autocompletion script for bash",
		Run: func(cmd *cobra.Command, args []string) {
			var filename, binaryName string
			switch runtime.GOOS {
			case "linux":
				filename = "/etc/bash_completion.d/instant-linux"
				binaryName = "instant-linux"

				err := os.Remove(filename)
				if err != nil && !os.IsNotExist(err) {
					log.Error(context.Background(), err)
					panic(err)
				}
			case "darwin":
				output, err := execShell.Exec("bash", "-c", "echo $(brew --prefix)")
				if err != nil {
					log.Error(context.Background(), err)
					panic(err)
				}
				filename = output + "/_instant-macos"
				binaryName = "instant-macos"

				err = os.Remove(filename)
				if err != nil && !os.IsNotExist(err) {
					log.Error(context.Background(), err)
					panic(err)
				}
			case "windows":
				log.Error(context.Background(), errors.New("autocomplete not supported for windows powershell"))
				panic("")
			}

			file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			err = cmd.GenBashCompletionV2(file, true)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			err = fileUtil.Sed(filename, "bash", binaryName)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			fmt.Println("Reload your shell session to begin using autocomplete!")
		},
	}

	return cmd
}
