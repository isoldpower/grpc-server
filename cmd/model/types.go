package model

import "github.com/spf13/cobra"

type SubCommand interface {
	Register(parentCmd *cobra.Command)
}

type Command interface {
	Execute() error
}
