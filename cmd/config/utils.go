package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

func resolveViper(viper *viper.Viper, configPath string) {
	configDir, configName, configType := splitViperPath(configPath)

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
}

func tryResolveConfig(viper *viper.Viper) error {
	err := viper.ReadInConfig()
	return err
}

func splitViperPath(path string) (configDir string, configName string, configType string) {
	configDir = filepath.Dir(path)
	configFullName := filepath.Base(path)
	configType = strings.TrimPrefix(filepath.Ext(configFullName), ".")
	configName = strings.TrimSuffix(configFullName, "."+configType)

	return configDir, configName, configType
}
