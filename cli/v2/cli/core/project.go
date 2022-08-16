package core

import (
	"io/ioutil"

	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"gopkg.in/yaml.v2"
)

func GenerateConfigFile(config *Config) {
	data, err := yaml.Marshal(&config)
	util.LogError(err)

	err = ioutil.WriteFile("config.yaml", data, 0)
	util.LogError(err)
}
