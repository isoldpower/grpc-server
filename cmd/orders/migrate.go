package orders

import (
	"github.com/spf13/cobra"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
)

type MigrateCommand struct {
	config          *Config
	commandInstance *cobra.Command
}

func NewMigrateCommand(config *Config) *MigrateCommand {
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
				log.Infoln("Executed migrate orders command")
				log.Debugln("Resolved orders config: %s", log.GetObjectPattern(config.Store))
			},
		},
	}
}

func (rc *MigrateCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
