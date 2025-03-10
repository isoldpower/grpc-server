package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
)

type CliConfigKey string

const (
	DebugCliConfigKey       CliConfigKey = "debug"
	SilentCliConfigKey      CliConfigKey = "silent"
	SkipClarifyCliConfigKey CliConfigKey = "skip-clarify"
	NoIconsCliConfigKey     CliConfigKey = "no-icons"
	NoColorsCliConfigKey    CliConfigKey = "no-colors"
)

type CliConfig struct {
	Silent      bool
	Debug       bool
	NoIcons     bool
	SkipClarify bool
	NoColors    bool

	viperInstance *viper.Viper
}

func (cc *CliConfig) applyValues() {
	log.SwitchDebug(cc.Debug)
	log.SwitchSilent(cc.Silent)
	log.SwitchIcons(!cc.NoIcons)

	color.SetEnabled(!cc.NoColors)
}

func NewCliConfig(viperInstance *viper.Viper) *CliConfig {
	cliConfig := &CliConfig{
		viperInstance: viperInstance,
		Silent:        false,
		Debug:         false,
		SkipClarify:   false,
		NoColors:      false,
	}

	return cliConfig
}

func (cc *CliConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&cc.Silent, string(SilentCliConfigKey), cc.Silent, "reduce the output to obligatory-only")
	cmd.PersistentFlags().BoolVar(&cc.Debug, string(DebugCliConfigKey), cc.Debug, "increase the amount of output information and print debug information")
	cmd.PersistentFlags().BoolVar(&cc.NoIcons, string(NoIconsCliConfigKey), cc.NoIcons, "exclude prefix icons from the logging")
	cmd.PersistentFlags().BoolVar(&cc.SkipClarify, string(SkipClarifyCliConfigKey), cc.SkipClarify, "automatically accept all 'yes/no' questions")
	cmd.PersistentFlags().BoolVar(&cc.NoColors, string(NoColorsCliConfigKey), cc.NoColors, "disable colorful output for CLI")
}

func (cc *CliConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver ParamReader = NewDualReader(cc.viperInstance, flags)

	cc.Debug = resolver.SafeGetBool(string(DebugCliConfigKey), cc.Debug)
	cc.Silent = resolver.SafeGetBool(string(SilentCliConfigKey), cc.Silent)
	cc.SkipClarify = resolver.SafeGetBool(string(SkipClarifyCliConfigKey), cc.SkipClarify)
	cc.NoIcons = resolver.SafeGetBool(string(NoIconsCliConfigKey), cc.NoIcons)
	cc.NoColors = resolver.SafeGetBool(string(NoColorsCliConfigKey), cc.NoColors)

	cc.applyValues()

	return nil
}

func (cc *CliConfig) TryReadConfig(path string) error {
	ResolveViper(cc.viperInstance, path)
	err := TryResolveConfig(cc.viperInstance)
	if err != nil {
		log.PrintError(fmt.Sprintf("Failed to resolve viper config at path: %s\n", path), err)
	}

	return nil
}
