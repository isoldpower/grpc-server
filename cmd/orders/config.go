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

	prefix        string
	serviceConfig string
	viperInstance *viper.Viper
}

func NewOrdersConfig(rootConfig *config.RootConfig) *Config {
	return &Config{
		Store: &types.InitialConfig{
			Root: rootConfig,
			Test: "default",
		},

		prefix:        "",
		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance: viper.New(),
	}
}

func NewPrefixedOrdersConfig(rootConfig *config.RootConfig, prefix string) *Config {
	return &Config{
		Store: &types.InitialConfig{
			Root: rootConfig,
			Test: "default",
		},

		prefix:        prefix,
		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance: viper.New(),
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
	cmd.PersistentFlags().StringVar(
		&oc.Store.Test,
		applier.WithPrefix(string(TestConfigKey)),
		oc.Store.Test,
		"just test variable",
	)
}

func (oc *Config) TryResolveConfig(_ string) error {
	config.ResolveViper(oc.viperInstance, oc.serviceConfig)
	err := config.TryResolveConfig(oc.viperInstance)
	if err != nil {
		return err
	}

	return nil
}

func (oc *Config) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver config.ParamReader = config.NewDualReader(oc.viperInstance, flags)

	oc.Store.Test = resolver.SafeGetString(string(TestConfigKey), oc.Store.Test)

	return nil
}
