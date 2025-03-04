package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ordersCommand *cobra.Command = &cobra.Command{
		Use:   "orders",
		Short: "Run Orders microservice",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//TODO: resolve orders-specific config options (i.e. database connection):
			// { runConfig: {}, ordersOptions...}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: pass resolved config to the function
			// orders.StartOrdersService()

			fmt.Println("Executed 'run orders' command. Running orders microservice")
		},
	}
)
