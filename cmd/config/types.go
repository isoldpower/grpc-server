package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandConfig interface {
	RegisterFlags(cmd *cobra.Command)
	TryResolveConfig(path string) error
	ResolveFlagsAndArgs(flagSet *pflag.FlagSet, args []string) error
}

type ParamReader interface {
	SafeGetString(key string, current string) string
	SafeGetBool(key string, current bool) bool
}
