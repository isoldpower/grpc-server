package orders

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/internal/util"
)

type RunCommand struct {
	config          *ordersConfig
	commandInstance *cobra.Command
}

func NewRunCommand(ordersConfig *ordersConfig) *RunCommand {
	return &RunCommand{
		config: ordersConfig,
		commandInstance: &cobra.Command{
			Use:   "run",
			Short: "Run Orders microservice",
			Long:  "Runs the whole process of Orders microservice",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Parent().PreRunE(cmd, args), func() error {
					return nil
				})
			},
			Run: func(cmd *cobra.Command, args []string) {
				value, _ := json.MarshalIndent(ordersConfig, "", "  ")
				fmt.Printf("Executed run orders command. Resolved config: %s\n", value)
			},
		},
	}
}

func (rc *RunCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
