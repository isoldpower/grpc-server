package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/model"
	"golang-grpc/cmd/orders"
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
				value, _ := json.MarshalIndent(rootConfig, "", "  ")
				fmt.Printf("Executed root command. Resolved config: %s\n", value)
			},
		},
	}
}

func init() {
	currentCommand.rootConfig.RegisterFlags(currentCommand.commandInstance)

	subcommands := []model.SubCommand{
		NewRunCommand(currentCommand.rootConfig),
		orders.NewRootCommand(currentCommand.rootConfig),
	}

	for _, subcommand := range subcommands {
		subcommand.Register(currentCommand.commandInstance)
	}
}

// Execute is an entry-point function to start the CLI interactions
func (c *RootCommand) Execute() error {
	if err := currentCommand.commandInstance.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}
