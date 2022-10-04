package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetConfigViper(configFile string) (*viper.Viper, error) {
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
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return configViper, nil
}

func GetEnvironmentVariableViper(envFiles []string) (*viper.Viper, error) {
	envVarViper := viper.New()

	if len(envFiles) > 0 {
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
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)

		envVarViper.AddConfigPath(wd)
		envVarViper.SetConfigType("env")
		envVarViper.SetConfigName(".env")

		_, err = os.Stat(filepath.Join(wd, ".env"))
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		err = envVarViper.ReadInConfig()
		if err != nil {
			return nil, errors.Wrap(err, "")
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
