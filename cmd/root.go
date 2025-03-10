package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/kitchen"
	"golang-grpc/cmd/orders"
	"golang-grpc/cmd/types"
	"golang-grpc/internal/log"
	"os"
)

type RootCommand struct {
	rootConfig      *config.RootConfig
	commandInstance *cobra.Command
}

var currentCommand = NewCommand()

func NewCommand() *RootCommand {
	rootConfig := config.NewRootConfig()

	return &RootCommand{
		rootConfig: rootConfig,
		commandInstance: &cobra.Command{
			Use:     "power",
			Version: "1.0.0",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return rootConfig.TryResolveConfig("")
			},
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return rootConfig.ResolveFlagsAndArgs(cmd.Flags(), args)
			},
			Run: func(cmd *cobra.Command, args []string) {
				log.Infoln("Executed root command")
				log.Debugln("Resolved config: %s\n", log.GetObjectPattern(rootConfig))
			},
		},
	}
}

func init() {
	currentCommand.rootConfig.RegisterFlags(currentCommand.commandInstance)

	subcommands := []types.SubCommand{
		NewRunCommand(currentCommand.rootConfig),
		orders.NewRootCommand(currentCommand.rootConfig),
		kitchen.NewRootCommand(currentCommand.rootConfig),
	}

	for _, subcommand := range subcommands {
		subcommand.Register(currentCommand.commandInstance)
	}
}

// Execute is an entry-point function to start the CLI interactions
func (c *RootCommand) Execute() error {
	if err := currentCommand.commandInstance.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		log.Errorln(err.Error())
		return err
	}

	return nil
}
