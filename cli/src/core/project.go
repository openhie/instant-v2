package core

import (
	"io/ioutil"

	"github.com/luno/jettison/errors"
	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidConfig = errors.New("invalid project config, required fields are Image, ProjectName, and PlatformImage")
)

func GenerateConfigFile(config *Config) error {
	if config.Image == "" || config.ProjectName == "" || config.PlatformImage == "" {
		return errors.Wrap(ErrInvalidConfig, "")
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = ioutil.WriteFile("config.yaml", data, 0600)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
