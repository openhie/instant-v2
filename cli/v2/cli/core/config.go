package core

import (
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/viper"
)

func LoadConfig(path string) Config {
	var config Config
	configViper := viper.New()
	configViper.SetConfigFile("config.yml")

	err := configViper.ReadInConfig()
	util.LogError(err)

	err = configViper.Unmarshal(&config)
	util.LogError(err)
	return config
}
