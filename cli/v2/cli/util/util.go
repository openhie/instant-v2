package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LogError(err error) {
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
