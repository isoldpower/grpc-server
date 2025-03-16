package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/internal/database"
	"golang-grpc/internal/log"
)

type DatabaseCliKey string
type DatabaseConfigKey string

const (
	HostCliKey     DatabaseCliKey = "db-host"
	PortCliKey     DatabaseCliKey = "db-port"
	UsernameCliKey DatabaseCliKey = "db-username"
	PasswordCliKey DatabaseCliKey = "db-password"
	NameCliKey     DatabaseCliKey = "db-name"
	SchemaCliKey   DatabaseCliKey = "db-schema"
)

const (
	HostConfigKey     DatabaseConfigKey = "database.host"
	PortConfigKey     DatabaseConfigKey = "database.port"
	UsernameConfigKey DatabaseConfigKey = "database.username"
	PasswordConfigKey DatabaseConfigKey = "database.password"
	NameConfigKey     DatabaseConfigKey = "database.name"
	SchemaConfigKey   DatabaseConfigKey = "database.schema"
)

type DatabaseConfig struct {
	*database.Config

	viperInstance *viper.Viper
}

func NewDatabaseConfig(viperInstance *viper.Viper) *DatabaseConfig {
	dbConfig := &DatabaseConfig{
		viperInstance: viperInstance,
		Config: &database.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "",
			Password: "",
			Database: "",
			Schema:   "public",
		},
	}

	return dbConfig
}

func (dc *DatabaseConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&dc.Config.Host, string(HostCliKey), dc.Config.Host, "set database host")
	cmd.Flags().StringVar(&dc.Config.Database, string(NameCliKey), dc.Config.Database, "set target database name")
	cmd.Flags().StringVar(&dc.Config.Username, string(UsernameCliKey), dc.Config.Username, "set database username")
	cmd.Flags().StringVar(&dc.Config.Password, string(PasswordCliKey), dc.Config.Password, "set database access password")
	cmd.Flags().StringVar(&dc.Config.Schema, string(SchemaCliKey), dc.Config.Schema, "set database search schema")
	cmd.Flags().IntVar(&dc.Config.Port, string(PortCliKey), dc.Config.Port, "set database port")
}

func (dc *DatabaseConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver ParamReader = NewDualReader(dc.viperInstance, flags)

	dc.Config.Host = resolver.SafeGetString(string(HostConfigKey), dc.Config.Host)
	dc.Config.Database = resolver.SafeGetString(string(NameConfigKey), dc.Config.Database)
	dc.Config.Username = resolver.SafeGetString(string(UsernameConfigKey), dc.Config.Username)
	dc.Config.Password = resolver.SafeGetString(string(PasswordConfigKey), dc.Config.Password)
	dc.Config.Schema = resolver.SafeGetString(string(SchemaConfigKey), dc.Config.Schema)
	dc.Config.Port = resolver.SafeGetInt(string(PortConfigKey), dc.Config.Port)

	return nil
}

func (dc *DatabaseConfig) TryReadConfig(path string) error {
	ResolveViper(dc.viperInstance, path)
	err := TryResolveConfig(dc.viperInstance)
	if err != nil {
		log.PrintError(fmt.Sprintf("Failed to resolve viper config at path: %s\n", path), err)
	}

	return nil
}
