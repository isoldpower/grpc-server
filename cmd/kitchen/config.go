package kitchen

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/cmd/config"
	"golang-grpc/services/kitchen/store"
	"path/filepath"
)

type configKey string

const (
	TestConfigKey configKey = "test"
	ConfigKey     configKey = "kitchen-config"
)

type KitchenConfig struct {
	store *store.InitialConfig

	serviceConfig string
	viperInstance *viper.Viper
}

func NewKitchenConfig(rootConfig *config.RootConfig) *KitchenConfig {
	return &KitchenConfig{
		store: &store.InitialConfig{
			Root: rootConfig,
			Test: "default",
		},

		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "kitchen", "config.yaml"),
		viperInstance: viper.New(),
	}
}

func (oc *KitchenConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&oc.serviceConfig, string(ConfigKey), oc.serviceConfig, "change service-specific config path")
	cmd.PersistentFlags().StringVar(&oc.store.Test, string(TestConfigKey), oc.store.Test, "")
}

func (oc *KitchenConfig) TryResolveConfig(_ string) error {
	config.ResolveViper(oc.viperInstance, oc.serviceConfig)
	err := config.TryResolveConfig(oc.viperInstance)
	if err != nil {
		return err
	}

	return nil
}

func (oc *KitchenConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver config.ParamReader = config.NewDualReader(oc.viperInstance, flags)

	oc.store.Test = resolver.SafeGetString(string(TestConfigKey), oc.store.Test)

	return nil
}
