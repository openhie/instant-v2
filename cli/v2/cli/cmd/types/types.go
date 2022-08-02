package types

import "github.com/spf13/viper"

type Global struct {
	ConfigViper *viper.Viper
	EnvVarViper *viper.Viper
}
