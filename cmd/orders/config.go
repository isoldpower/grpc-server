package orders

import (
	"fmt"
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
	grpcConfig     *config.ServerConfig
	httpConfig     *config.ServerConfig
	viperInstance  *viper.Viper
}

func NewOrdersConfig(rootConfig *config.RootConfig) *Config {
	viperInstance := viper.New()
	databaseConfig := config.NewDatabaseConfig(viperInstance)
	grpcConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 3081, "grpc")
	httpConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 3082, "http")

	return &Config{
		Store: &types.InitialConfig{
			Root:     rootConfig,
			Database: databaseConfig.Config,
			GRPC:     grpcConfig.ServerConfig,
			HTTP:     httpConfig.ServerConfig,
		},

		prefix:         "",
		serviceConfig:  filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance:  viperInstance,
		databaseConfig: databaseConfig,
		grpcConfig:     grpcConfig,
		httpConfig:     httpConfig,
	}
}

func NewPrefixedOrdersConfig(rootConfig *config.RootConfig, prefix string) *Config {
	viperInstance := viper.New()
	databaseConfig := config.NewDatabaseConfig(viperInstance)
	grpcConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 3081, fmt.Sprintf("%s-grpc", prefix))
	httpConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 3082, fmt.Sprintf("%s-http", prefix))

	return &Config{
		Store: &types.InitialConfig{
			Root:     rootConfig,
			Database: databaseConfig.Config,
			GRPC:     grpcConfig.ServerConfig,
			HTTP:     httpConfig.ServerConfig,
		},

		prefix:         prefix,
		serviceConfig:  filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance:  viperInstance,
		databaseConfig: databaseConfig,
		grpcConfig:     grpcConfig,
		httpConfig:     httpConfig,
	}
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
	oc.grpcConfig.RegisterFlags(cmd)
	oc.httpConfig.RegisterFlags(cmd)
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
	if err := oc.grpcConfig.ResolveFlagsAndArgs(flags, args); err != nil {
		return err
	}
	if err := oc.httpConfig.ResolveFlagsAndArgs(flags, args); err != nil {
		return err
	}

	return nil
}
