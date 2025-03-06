package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
)

type RunCommand struct {
	commandInstance *cobra.Command
}

func NewRunCommand(rootConfig *config.RootConfig) *RunCommand {
	return &RunCommand{
		commandInstance: &cobra.Command{
			Use:   "run",
			Short: "Run all microservices in correct order",
			PreRunE: func(cmd *cobra.Command, args []string) error {
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
