package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.Reset()
	viper.SetConfigName("chlog")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if path != "" {
		viper.SetConfigFile(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if path != "" {
				return fmt.Errorf("Config file not found: %s", path)
			}
			return nil
		}
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func GetConfigFlagString(cmd *cobra.Command, name string) (string, bool, error) {
	flagValue, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", false, err
	}

	if flagValue != "" {
		return flagValue, false, nil
	}

	if viper.IsSet(name) {
		return viper.GetString(name), true, nil
	}

	return "", false, nil
}

func GetConfigFlagBool(cmd *cobra.Command, name string) (bool, error) {
	flagValue, err := cmd.Flags().GetBool(name)
	if err != nil {
		return false, err
	}

	if flagValue {
		return true, nil
	}

	if viper.IsSet(name) {
		return viper.GetBool(name), nil
	}


	return false, nil
}
