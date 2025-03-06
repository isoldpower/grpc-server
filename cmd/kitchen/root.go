package kitchen

import (
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
)

var (
	rootCommand = &cobra.Command{
		Use:   "kitchen",
		Short: "Kitchen Microservice CLI",
		Long:  `Kitchen Microservice CLI. Root command outputs help`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

type Command struct {
	rootConfig *config.RootConfig
	config     *KitchenConfig
}

func NewKitchenCommand(rootConfig *config.RootConfig) *Command {
	return &Command{
		rootConfig: rootConfig,
	}
}

func (c *Command) Register(parentCmd *cobra.Command) {

}
