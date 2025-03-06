package orders

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/model"
	"golang-grpc/internal/util"
)

type RootCommand struct {
	ordersConfig    *ordersConfig
	commandInstance *cobra.Command
}

func NewRootCommand(rootConfig *config.RootConfig) *RootCommand {
	commandConfig := newOrdersConfig(rootConfig)

	return &RootCommand{
		ordersConfig: commandConfig,
		commandInstance: &cobra.Command{
			Use:   "orders",
			Short: "Orders microservice",
			Long:  "Orders microservice-specific CLI for controlling Orders microservice.\nRoot command outputs help by default",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Root().PersistentPreRunE(cmd, args), func() error {
					return commandConfig.TryResolveConfig("")
				})
			},
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Root().PreRunE(cmd, args), func() error {
					return commandConfig.ResolveFlagsAndArgs(cmd.Flags(), args)
				})
			},
			Run: func(cmd *cobra.Command, args []string) {
				value, _ := json.MarshalIndent(commandConfig, "", "  ")
				fmt.Printf("Executed orders command. Resolved config: %s\n", value)
			},
		},
	}
}

func (rc *RootCommand) Register(parentCmd *cobra.Command) {
	rc.ordersConfig.RegisterFlags(rc.commandInstance)

	subCommands := []model.SubCommand{
		NewRunCommand(rc.ordersConfig),
		NewMigrateCommand(rc.ordersConfig),
	}

	for _, subCommand := range subCommands {
		subCommand.Register(rc.commandInstance)
	}

	parentCmd.AddCommand(rc.commandInstance)
}
