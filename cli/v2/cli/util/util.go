package util

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Log(message string) {
	if os.Getenv("LOG") == "true" {
		fmt.Println(message)
	}
}

func LogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func PanicError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func GetFlagOrDefaultString(cmd *cobra.Command, flagName string) string {
	var name string
	if cmd.Flag(flagName).Changed {
		var err error
		name, err = cmd.Flags().GetString(flagName)
		LogError(err)
	} else {
		name = cmd.Flag(flagName).DefValue
	}
	return name

}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
