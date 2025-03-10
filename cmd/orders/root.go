package orders

import (
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/types"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
)

type RootCommand struct {
	ordersConfig    *Config
	commandInstance *cobra.Command
}

func NewRootCommand(rootConfig *config.RootConfig) *RootCommand {
	commandConfig := NewOrdersConfig(rootConfig)

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
				log.Infoln("Executed root orders command")
				log.Debugln("Resolved orders config: %s", log.GetObjectPattern(commandConfig.Store))
				cmd.HelpFunc()(cmd, args)
			},
		},
	}
}

func (rc *RootCommand) Register(parentCmd *cobra.Command) {
	rc.ordersConfig.RegisterFlags(rc.commandInstance)

	subCommands := []types.SubCommand{
		NewRunCommand(rc.ordersConfig),
		NewMigrateCommand(rc.ordersConfig),
	}

	for _, subCommand := range subCommands {
		subCommand.Register(rc.commandInstance)
	}

	parentCmd.AddCommand(rc.commandInstance)
}
