package orders

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/internal/util"
)

type MigrateCommand struct {
	config          *ordersConfig
	commandInstance *cobra.Command
}

func NewMigrateCommand(config *ordersConfig) *MigrateCommand {
	return &MigrateCommand{
		config: config,
		commandInstance: &cobra.Command{
			Use:   "migrate",
			Short: "Migrate Orders microservice",
			Long:  "Run further CLI process of migrating all Orders services",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Parent().PreRunE(cmd, args), func() error {
					return nil
				})
			},
			Run: func(cmd *cobra.Command, args []string) {
				value, _ := json.MarshalIndent(config, "", "  ")
				fmt.Printf("Executed migrate orders command. Resolved config: %s\n", value)
			},
		},
	}
}

func (rc *MigrateCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
