package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	kitchenCommand *cobra.Command = &cobra.Command{
		Use:   "kitchen",
		Short: "Run Kitchen microservice",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//TODO: resolve kitchen-specific config options (i.e. gRPC orders connection):
			// { runConfig: {}, kitchenOptions...}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: pass resolved config to the function
			// kitchen.StartKitchenService()

			fmt.Println("Executed 'run kitchen' command. Running kitchen microservice")
		},
	}
)
