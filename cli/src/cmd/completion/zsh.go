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

func genZshCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zsh",
		Short: "Generate the autocompletion script for zsh",
		Run: func(cmd *cobra.Command, args []string) {
			var filename, binaryName string
			switch runtime.GOOS {
			case "linux":
				output, err := execShell.Exec("zsh", "-c", "echo ${fpath[1]}")
				if err != nil {
					log.Error(context.Background(), err)
					panic(err)
				}
				filename = output + "/_gocli-linux"
				binaryName = "gocli-linux"

				err = os.Remove(filename)
				if err != nil && !os.IsNotExist(err) {
					log.Error(context.Background(), err)
					panic(err)
				}
			case "darwin":
				output, err := execShell.Exec("zsh", "-c", "echo $(brew --prefix)")
				if err != nil {
					log.Error(context.Background(), err)
					panic(err)
				}
				filename = output + "/_gocli-macos"
				binaryName = "gocli-macos"

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

			err = cmd.GenZshCompletion(file)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			err = fileUtil.Sed(filename, "zsh", binaryName)
			if err != nil {
				log.Error(context.Background(), err)
				panic(err)
			}

			fmt.Println("Reload your shell session to begin using autocomplete!")
		},
	}

	return cmd
}
