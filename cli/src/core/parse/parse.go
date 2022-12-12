package parse

import (
	"strings"

	"cli/core"
	coreConfig "cli/core/state"

	"github.com/luno/jettison/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidConfigFileSyntax = errors.New("invalid config file syntax, refer to https://github.com/openhie/package-starter-kit/blob/main/README.md, for information on valid config file syntax")
)

func unmarshalConfig(configViper *viper.Viper) (*core.Config, error) {
	var config core.Config
	err := configViper.Unmarshal(&config)
	if err != nil && strings.Contains(err.Error(), "expected type") {
		return nil, errors.Wrap(ErrInvalidConfigFileSyntax, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &config, nil
}

func appendTag(config *core.Config) {
	splitStrings := strings.Split(config.Image, ":")

	if len(splitStrings) == 1 {
		config.Image += ":latest"
	}
}

func getConfigFromParams(cmd *cobra.Command) (*core.Config, error) {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	configViper, err := coreConfig.SetConfigViper(configFile)
	if err != nil {
		return nil, err
	}

	populatedConfig, err := unmarshalConfig(configViper)
	if err != nil {
		return nil, err
	}

	appendTag(populatedConfig)

	return populatedConfig, nil
}
