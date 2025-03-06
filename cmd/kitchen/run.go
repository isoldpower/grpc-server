package kitchen

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/internal/util"
	"golang-grpc/services/kitchen"
)

type RunCommand struct {
	config          *KitchenConfig
	commandInstance *cobra.Command
}

func NewRunCommand(kitchenConfig *KitchenConfig) *RunCommand {
	return &RunCommand{
		config: kitchenConfig,
		commandInstance: &cobra.Command{
			Use:   "run",
			Short: "Run Kitchen microservice",
			Long:  "Runs the whole process of Kitchen microservice.\n\tPay attention that Kitchen service is dependant on Orders service",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Parent().PreRunE(cmd, args), func() error {
					return nil
				})
			},
			Run: func(cmd *cobra.Command, args []string) {
				value, _ := json.MarshalIndent(kitchenConfig.store, "", "  ")
				fmt.Printf("Executed run kitchen command. Resolved config: %s\n", value)

				kitchen.StartKitchenService(kitchenConfig.store)
			},
		},
	}
}

func (rc *RunCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
