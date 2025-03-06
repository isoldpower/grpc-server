package orders

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/cmd/config"
	"path/filepath"
)

type configKey string

const (
	TestConfigKey configKey = "test"
	ConfigKey     configKey = "orders-config"
)

type ordersConfig struct {
	Root *config.RootConfig
	Test string

	serviceConfig string
	viperInstance *viper.Viper
}

func newOrdersConfig(rootConfig *config.RootConfig) *ordersConfig {
	return &ordersConfig{
		Root: rootConfig,
		Test: "default",

		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "orders", "config.yaml"),
		viperInstance: viper.New(),
	}
}

func (oc *ordersConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&oc.serviceConfig, string(ConfigKey), oc.serviceConfig, "change service-specific config path")
}

func (oc *ordersConfig) TryResolveConfig(_ string) error {
	config.ResolveViper(oc.viperInstance, oc.serviceConfig)
	err := config.TryResolveConfig(oc.viperInstance)
	if err != nil {
		return err
	}

	return nil
}

func (oc *ordersConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver config.ParamReader = config.NewDualReader(oc.viperInstance, flags)
	oc.Test = resolver.SafeGetString(string(TestConfigKey), oc.Test)

	return nil
}
