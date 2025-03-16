package orders

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/cmd/config"
	"golang-grpc/internal/util"
	"golang-grpc/services/orders/types"
	"path/filepath"
)

type configKey string

const (
	TestConfigKey configKey = "test"
	ConfigKey     configKey = "orders-config"
)

type Config struct {
	Store *types.InitialConfig

	prefix         string
	serviceConfig  string
	databaseConfig *config.DatabaseConfig
	viperInstance  *viper.Viper
}

func NewOrdersConfig(rootConfig *config.RootConfig) *Config {
	viperInstance := viper.New()
	databaseConfig := config.NewDatabaseConfig(viperInstance)

	return &Config{
		Store: &types.InitialConfig{
			Root:     rootConfig,
			Database: databaseConfig.Config,
		},

		prefix:         "",
		serviceConfig:  filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance:  viperInstance,
		databaseConfig: databaseConfig,
	}
}

func NewPrefixedOrdersConfig(rootConfig *config.RootConfig, prefix string) *Config {
	configInstance := NewOrdersConfig(rootConfig)
	configInstance.prefix = prefix

	return configInstance
}

func (oc *Config) RegisterFlags(cmd *cobra.Command) {
	applier := util.NewPrefixApplier(oc.prefix)

	cmd.PersistentFlags().StringVar(
		&oc.serviceConfig,
		applier.WithPrefix(string(ConfigKey)),
		oc.serviceConfig,
		"change service-specific config path",
	)
	oc.databaseConfig.RegisterFlags(cmd)
}

func (oc *Config) TryResolveConfig(_ string) error {
	config.ResolveViper(oc.viperInstance, oc.serviceConfig)
	err := config.TryResolveConfig(oc.viperInstance)
	if err != nil {
		return err
	}

	return nil
}

func (oc *Config) ResolveFlagsAndArgs(flags *pflag.FlagSet, args []string) error {
	if err := oc.databaseConfig.ResolveFlagsAndArgs(flags, args); err != nil {
		return err
	}

	return nil
}
