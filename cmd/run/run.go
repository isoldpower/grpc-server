package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	runCommand *cobra.Command = &cobra.Command{
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
	}
)

func init() {
	runCommand.AddCommand(ordersCommand)
	runCommand.AddCommand(kitchenCommand)
}

// GetCommand returns reference to runCommand instance
func GetCommand() *cobra.Command {
	return runCommand
}
