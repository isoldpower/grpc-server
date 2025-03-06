package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type RootConfigKey string

type RootConfig struct {
	Cli           *CliConfig
	Context       *ProcessContext
	configPath    string
	viperInstance *viper.Viper
}

func NewRootConfig() *RootConfig {
	rootViper := viper.New()

	return &RootConfig{
		viperInstance: rootViper,
		Cli:           NewCliConfig(rootViper),
		Context:       NewProcessContext(),
	}
}

func (cc *RootConfig) RegisterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cc.configPath, "config", "./root_config.yaml", "config file path")

	cc.Cli.RegisterFlags(cmd)
}

func (cc *RootConfig) ResolveFlagsAndArgs(flagSet *pflag.FlagSet, args []string) error {
	contextError := cc.Context.ResolveProcessContext(args)
	if contextError != nil {
		fmt.Printf("Failed to resolve process context: %v\n", contextError)
		return contextError
	}

	cliError := cc.Cli.ResolveFlagsAndArgs(flagSet, args)
	if cliError != nil {
		fmt.Printf("Failed to resolve CLI args and flags: %v\n", cliError)
		return cliError
	}

	return nil
}

func (cc *RootConfig) TryResolveConfig(_ string) error {
	cliError := cc.Cli.TryReadConfig(cc.configPath)
	if cliError != nil {
		fmt.Printf("Failed to read config: %v\n", cliError)
		return cliError
	}

	return nil
}
