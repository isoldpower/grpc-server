package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CliConfigKey string

const (
	DebugCliConfigKey       CliConfigKey = "debug"
	SilentCliConfigKey      CliConfigKey = "silent"
	SkipClarifyCliConfigKey CliConfigKey = "skip-clarify"
)

type CliConfig struct {
	Silent      bool
	Debug       bool
	SkipClarify bool

	viperInstance *viper.Viper
}

func NewCliConfig(viperInstance *viper.Viper) *CliConfig {
	cliConfig := &CliConfig{
		viperInstance: viperInstance,
		Silent:        false,
		Debug:         false,
		SkipClarify:   false,
	}

	return cliConfig
}

func (cc *CliConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&cc.Silent, string(SilentCliConfigKey), cc.Silent, "reduce the output to obligatory-only")
	cmd.PersistentFlags().BoolVar(&cc.Debug, string(DebugCliConfigKey), cc.Debug, "increase the amount of output information and print debug information")
	cmd.PersistentFlags().BoolVar(&cc.SkipClarify, string(SkipClarifyCliConfigKey), cc.SkipClarify, "automatically accept all 'yes/no' questions")
}

func (cc *CliConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver ParamReader = NewDualReader(cc.viperInstance, flags)

	cc.Debug = resolver.SafeGetBool(string(DebugCliConfigKey), cc.Debug)
	cc.Silent = resolver.SafeGetBool(string(SilentCliConfigKey), cc.Silent)
	cc.Debug = resolver.SafeGetBool(string(SkipClarifyCliConfigKey), cc.SkipClarify)

	return nil
}

func (cc *CliConfig) TryReadConfig(path string) error {
	ResolveViper(cc.viperInstance, path)
	err := TryResolveConfig(cc.viperInstance)
	if err != nil {
		fmt.Printf("Failed to resolve viper config at path: %s\n", path)
		fmt.Println(err.Error())
	}

	return nil
}
