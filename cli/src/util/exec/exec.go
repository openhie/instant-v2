package exec

import (
	"bufio"
	"os/exec"

	"github.com/pkg/errors"
)

func Exec(commandName string, commandSlice ...string) (string, error) {
	cmd := exec.Command(commandName, commandSlice...)
	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "Error creating stdOutPipe for Cmd.")
	}
	stdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.Wrap(err, "Error creating stdErrPipe for Cmd.")
	}

	var output []rune
	stdOutScanner := bufio.NewScanner(stdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			s := stdOutScanner.Text()
			if s != "" {
				for _, ss := range s {
					output = append(output, ss)
				}
			}
		}
	}()

	var errStr []rune
	stdErrScanner := bufio.NewScanner(stdErrReader)
	go func() {
		for stdErrScanner.Scan() {
			s := stdErrScanner.Text()
			if s != "" {
				for _, ss := range s {
					errStr = append(errStr, ss)
				}
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, string(errStr))
	}

	errString := string(errStr)
	if errString != "" {
		return "", errors.Wrap(errors.New(errString), "")
	}

	return string(output), nil
}
