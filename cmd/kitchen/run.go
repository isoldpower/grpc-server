package kitchen

import (
	"github.com/spf13/cobra"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
	"golang-grpc/services/kitchen"
)

type RunCommand struct {
	config          *Config
	commandInstance *cobra.Command
}

func NewRunCommand(kitchenConfig *Config) *RunCommand {
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
				log.Infoln("Executed migrate kitchen command")
				log.Debugln("Resolved kitchen config: %s", log.GetObjectPattern(kitchenConfig.Store))

				ready := make(chan bool, 1)
				service := kitchen.NewKitchenService(kitchenConfig.Store)
				done := service.Execute(ready)
				<-done
			},
		},
	}
}

func (rc *RunCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
