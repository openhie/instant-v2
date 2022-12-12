package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/spf13/viper"
)

var (
	ConfigFile  string
	EnvFiles    []string
	configViper *viper.Viper
)

func SetConfigViper(configFile string) (*viper.Viper, error) {
	configViper = viper.New()
	if configFile != "" {
		absFilePath, err := filepath.Abs(configFile)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		configViper.SetConfigFile(absFilePath)
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
		if !filepath.IsAbs(envFile) {
			envFile = filepath.Join(filepath.Dir(configViper.ConfigFileUsed()), envFile)
		}

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
