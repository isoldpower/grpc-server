package kitchen

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/cmd/config"
	"golang-grpc/internal/util"
	"golang-grpc/services/kitchen/store"
	"path/filepath"
)

type configKey string

const (
	TestConfigKey configKey = "test"
	ConfigKey     configKey = "kitchen-config"
)

type Config struct {
	Store *store.InitialConfig

	prefix        string
	serviceConfig string
	serverConfig  *config.ServerConfig
	viperInstance *viper.Viper
}

func NewKitchenConfig(rootConfig *config.RootConfig) *Config {
	viperInstance := viper.New()
	serverConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 8000, "")

	return &Config{
		Store: &store.InitialConfig{
			Root:   rootConfig,
			Server: serverConfig.ServerConfig,
			Test:   "default",
		},

		prefix:        "",
		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "kitchen", "config.yaml"),
		serverConfig:  serverConfig,
		viperInstance: viperInstance,
	}
}

func NewPrefixedKitchenConfig(rootConfig *config.RootConfig, prefix string) *Config {
	viperInstance := viper.New()
	serverConfig := config.NewServerConfig(viperInstance, "0.0.0.0", 8000, prefix)

	return &Config{
		Store: &store.InitialConfig{
			Root:   rootConfig,
			Server: serverConfig.ServerConfig,
			Test:   "default",
		},

		prefix:        prefix,
		serviceConfig: filepath.Join(rootConfig.Context.RootDir, "services", "kitchen", "config.yaml"),
		serverConfig:  serverConfig,
		viperInstance: viperInstance,
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
	oc.serverConfig.RegisterFlags(cmd)
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
	var resolver config.ParamReader = config.NewDualReader(oc.viperInstance, flags)

	oc.Store.Test = resolver.SafeGetString(string(TestConfigKey), oc.Store.Test)
	if err := oc.serverConfig.ResolveFlagsAndArgs(flags, args); err != nil {
		return err
	}

	return nil
}
