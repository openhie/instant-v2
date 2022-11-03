package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/spf13/viper"
)

var configViper = viper.New()

func SetConfigViper(configFile string) (*viper.Viper, error) {
	configViper = viper.New()
	if configFile != "" {
		configViper.SetConfigFile(configFile)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		configViper.AddConfigPath(wd)
		configViper.SetConfigType("yaml")
		configViper.SetConfigName("config")
	}

	err := configViper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return configViper, nil
}

func GetEnvironmentVariableViper(envFiles []string) (*viper.Viper, error) {
	envVarViper := viper.New()

	for i, envFile := range envFiles {
		_, err := os.Stat(envFile)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		envVarViper.SetConfigType("env")
		envVarViper.SetConfigFile(envFile)
		if i == 0 {
			err = envVarViper.ReadInConfig()
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
		} else {
			err := envVarViper.MergeInConfig()
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
		}
	}

	return envVarViper, nil
}

func GetEnvVariableString(envViper *viper.Viper) []string {
	var envVariables []string
	allEnvVars := envViper.AllSettings()
	for key, element := range allEnvVars {
		envVariables = append(envVariables, fmt.Sprintf("%v=%v", strings.ToUpper(key), element))
	}

	return envVariables
}
