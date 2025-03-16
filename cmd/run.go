package cmd

import (
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	kitchenCmd "golang-grpc/cmd/kitchen"
	ordersCmd "golang-grpc/cmd/orders"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
	"golang-grpc/services/common/types"
	"golang-grpc/services/kitchen"
	"golang-grpc/services/orders"
)

type RunCommand struct {
	commandInstance *cobra.Command
	services        map[string]types.Service
	kitchenConfig   *kitchenCmd.Config
	ordersConfig    *ordersCmd.Config
}

func (rc *RunCommand) runServicesInOrder(globalDone chan bool) {
	ready := make(chan bool)
	defer close(ready)
	doneChannels := make([]<-chan bool, len(rc.services))

	iterator := 0
	for key, service := range rc.services {
		go func() {
			log.Infoln("Running %s service", key)
			log.IncreaseLevel()
			doneChannels[iterator] = service.Execute(ready)
		}()
		<-ready
		log.DecreaseLevel()
		iterator++
	}
	log.Logln("\n")

	finalStream := util.FlatStreams(globalDone, doneChannels...)
	select {
	case globalDone <- <-finalStream:
		break
	}
}

func NewRunCommand(rootConfig *config.RootConfig) *RunCommand {
	kitchenConfig := kitchenCmd.NewPrefixedKitchenConfig(rootConfig, "kitchen")
	ordersConfig := ordersCmd.NewPrefixedOrdersConfig(rootConfig, "orders")
	configs := []config.CommandConfig{
		kitchenConfig, ordersConfig,
	}

	command := &RunCommand{
		kitchenConfig: kitchenConfig,
		ordersConfig:  ordersConfig,
		services: map[string]types.Service{
			"orders":  orders.NewOrdersService(ordersConfig.Store),
			"kitchen": kitchen.NewKitchenService(kitchenConfig.Store),
		},
	}

	command.commandInstance = &cobra.Command{
		Use:   "run",
		Short: "Run all microservices in correct order",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return util.ProtectedAction(cmd.Root().PersistentPreRunE(cmd, args), func() error {
				for _, cmdConfig := range configs {
					if err := cmdConfig.TryResolveConfig(""); err != nil {
						return err
					}
				}

				return nil
			})
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return util.ProtectedAction(cmd.Root().PreRunE(cmd, args), func() error {
				for _, cmdConfig := range configs {
					if err := cmdConfig.ResolveFlagsAndArgs(cmd.Flags(), args); err != nil {
						return err
					}
				}

				return nil
			})
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Infoln("Executed %s command", color.Underline("run-all"))
			log.Debugln("Resolved kitchen config: %s", log.GetObjectPattern(kitchenConfig.Store))
			log.Debugln("Resolved orders config: %s", log.GetObjectPattern(ordersConfig.Store))

			done := make(chan bool)
			command.runServicesInOrder(done)
		},
	}

	return command
}

func (rc *RunCommand) Register(parentCmd *cobra.Command) {
	rc.kitchenConfig.RegisterFlags(rc.commandInstance)
	rc.ordersConfig.RegisterFlags(rc.commandInstance)

	parentCmd.AddCommand(rc.commandInstance)
}
