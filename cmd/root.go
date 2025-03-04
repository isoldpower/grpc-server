package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/run"
	"os"
)

var (
	rootCommand = &cobra.Command{
		Use:     "power",
		Version: "1.0.0",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//TODO: try resolve running context and running root config (i.e. Logger configuration)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: output help for the commands
			fmt.Println("Executed root command")
		},
	}
)

func init() {
	rootCommand.AddCommand(run.GetCommand())
}

// Execute as an entry-point function to start the CLI interactions
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
