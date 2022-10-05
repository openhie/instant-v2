package core

import (
	"io/ioutil"

	"github.com/luno/jettison/errors"
	"gopkg.in/yaml.v2"
)

func GenerateConfigFile(config *Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = ioutil.WriteFile("config.yaml", data, 0)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
