package state

import (
	"os"
	"path/filepath"

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
