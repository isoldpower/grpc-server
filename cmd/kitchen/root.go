package kitchen

import (
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/types"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
)

type RootCommand struct {
	kitchenConfig   *Config
	commandInstance *cobra.Command
}

func NewRootCommand(rootConfig *config.RootConfig) *RootCommand {
	commandConfig := NewKitchenConfig(rootConfig)

	return &RootCommand{
		kitchenConfig: commandConfig,
		commandInstance: &cobra.Command{
			Use:   "kitchen",
			Short: "Kitchen microservice",
			Long:  "Kitchen microservice-specific CLI for controlling Kitchen microservice.\nRoot command outputs help by default",
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
				log.Infoln("Executed root kitchen command")
				log.Debugln("Resolved kitchen config: %s", log.GetObjectPattern(commandConfig.Store))
			},
		},
	}
}

func (rc *RootCommand) Register(parentCmd *cobra.Command) {
	rc.kitchenConfig.RegisterFlags(rc.commandInstance)

	subCommands := []types.SubCommand{
		NewRunCommand(rc.kitchenConfig),
		NewMigrateCommand(rc.kitchenConfig),
	}

	for _, subCommand := range subCommands {
		subCommand.Register(rc.commandInstance)
	}

	parentCmd.AddCommand(rc.commandInstance)
}
