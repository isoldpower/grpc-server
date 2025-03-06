package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
)

type RunCommand struct {
	rootConfig      *config.RootConfig
	commandInstance *cobra.Command
}

func NewRunCommand(rootConfig *config.RootConfig) *RunCommand {
	return &RunCommand{
		rootConfig: rootConfig,
		commandInstance: &cobra.Command{
			Use:   "run",
			Short: "Run the microservice",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				//TODO: resolve the configuration here
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {
				//TODO: run all microservices in correct integration order
				// in separate goroutines
				fmt.Println("Executed run command (same as 'run all')")
			},
		},
	}
}

func (rc *RunCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(rc.commandInstance)
}
