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

func NewCommand() *RootCommand {
	rootConfig := config.NewRootConfig()

	rootCommand := &RootCommand{
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
				cmd.HelpFunc()(cmd, args)
			},
		},
	}

	rootCommand.rootConfig.RegisterFlags(rootCommand.commandInstance)

	subcommands := []types.SubCommand{
		NewRunCommand(rootCommand.rootConfig),
		orders.NewRootCommand(rootCommand.rootConfig),
		kitchen.NewRootCommand(rootCommand.rootConfig),
	}

	for _, subcommand := range subcommands {
		subcommand.Register(rootCommand.commandInstance)
	}

	return rootCommand
}

// Execute is an entry-point function to start the CLI interactions
func (c *RootCommand) Execute() error {
	if err := c.commandInstance.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Errorln(err.Error())
		return err
	}

	return nil
}
