package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func ResolveViper(viper *viper.Viper, configPath string) {
	configDir, configName, configType := SplitViperPath(configPath)

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
}

func TryResolveConfig(viperInstance *viper.Viper) error {
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperInstance.AutomaticEnv()

	err := viperInstance.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func SplitViperPath(path string) (configDir string, configName string, configType string) {
	configDir = filepath.Dir(path)
	configFullName := filepath.Base(path)
	configType = strings.TrimPrefix(filepath.Ext(configFullName), ".")
	configName = strings.TrimSuffix(configFullName, "."+configType)

	return configDir, configName, configType
}
