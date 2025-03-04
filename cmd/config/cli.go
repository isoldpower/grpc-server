package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CliConfigKey string

const (
	DebugCliConfigKey  CliConfigKey = "debug"
	SilentCliConfigKey CliConfigKey = "silent"
)

type CliConfig struct {
	Silent        bool
	Debug         bool
	viperInstance *viper.Viper
}

func NewCliConfig(viperInstance *viper.Viper) *CliConfig {
	cliConfig := &CliConfig{
		viperInstance: viperInstance,
		Silent:        false,
		Debug:         false,
	}

	viperInstance.SetDefault(string(DebugCliConfigKey), cliConfig.Debug)
	viperInstance.SetDefault(string(SilentCliConfigKey), cliConfig.Silent)
	return cliConfig
}

func (cc *CliConfig) AttachFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&cc.Silent, string(SilentCliConfigKey), cc.Silent, "reduce the output to obligatory-only")
	cmd.Flags().BoolVar(&cc.Debug, string(DebugCliConfigKey), cc.Debug, "increase the amount of output information and print debug information")
}

func (cc *CliConfig) ResolveArgsAndFlags(flagSet *pflag.FlagSet, args []string) error {
	flagReader := NewFlagReader(flagSet)
	cc.Silent = flagReader.SafeGetBool(string(SilentCliConfigKey), cc.Silent)
	cc.Debug = flagReader.SafeGetBool(string(DebugCliConfigKey), cc.Debug)

	return nil
}

func (cc *CliConfig) TryReadConfig(path string) error {
	resolveViper(cc.viperInstance, path)
	err := tryResolveConfig(cc.viperInstance)
	if err != nil {
		fmt.Printf("Failed to resolve viper config at path: %s\n", path)
		fmt.Println(err.Error())
	}

	return nil
}
