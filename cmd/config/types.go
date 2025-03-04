package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandConfig interface {
	AttachFlagsToCommand(cmd *cobra.Command)
	TryReadConfig(path string) error
	ResolveArgsAndFlags(flags *pflag.FlagSet, args []string) error
}
