package cmd

import (
	"os"

	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/commands"
	"github.com/openhie/package-starter-kit/cli/v2/cli/cmd/types"
	"github.com/openhie/package-starter-kit/cli/v2/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var configViper viper.Viper

var envFiles []string
var envVarViper viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A cli to assist with package deployment and management",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initEnvironmentVariables)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $WORKING_DIR/config.yaml)")
	rootCmd.PersistentFlags().StringSliceVarP(&envFiles, "env-file", "e", nil, "env file (default is $WORKING_DIR/.env)")

	global := &types.Global{
		ConfigViper: &configViper,
		EnvVarViper: &envVarViper,
	}

	commands.AddCommands(rootCmd, global)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
	util.LogError(err)
}

func initEnvironmentVariables() {
	envVarViper := viper.New()

	if envFiles != nil {
		for i, envFile := range envFiles {
			envVarViper.SetConfigType("env")
			envVarViper.SetConfigFile(envFile)
			if i == 0 {
				err := envVarViper.ReadInConfig()
				util.LogError(err)
			} else {
				err := envVarViper.MergeInConfig()
				util.LogError(err)
			}
		}
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		envVarViper.AddConfigPath(wd)
		envVarViper.SetConfigType("env")
		envVarViper.SetConfigName(".env")
		err = envVarViper.ReadInConfig()
		util.LogError(err)
	}
}
