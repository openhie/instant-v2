package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetConfigViper(configFile string) *viper.Viper {
	configViper := viper.New()
	if configFile != "" {
		configViper.SetConfigFile(configFile)
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		configViper.AddConfigPath(wd)
		configViper.SetConfigType("yaml")
		configViper.SetConfigName("config")
	}

	err := configViper.ReadInConfig()
	util.LogError(err)
	return configViper
}

func GetEnvironmentVariableViper(envFiles []string) *viper.Viper {
	envVarViper := viper.New()

	if envFiles != nil {
		for i, envFile := range envFiles {
			envVarViper.SetConfigType("env")
			envVarViper.SetConfigFile(envFile)
			if i == 0 {
				err := envVarViper.ReadInConfig()
				util.LogError(err)
			} else {
				err := envVarViper.MergeInConfig()
				util.LogError(err)
			}
		}
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		envVarViper.AddConfigPath(wd)
		envVarViper.SetConfigType("env")
		envVarViper.SetConfigName(".env")
		err = envVarViper.ReadInConfig()
		util.LogError(err)
	}
	return envVarViper
}

func GetEnvVariableString(envViper *viper.Viper) []string {
	var envVariables []string
	allEnvVars := envViper.AllSettings()
	for key, element := range allEnvVars {
		envVariables = append(envVariables, fmt.Sprintf("%v=%v", strings.ToUpper(key), element))
	}
	return envVariables
}
