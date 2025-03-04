package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/cmd/config"
	"golang-grpc/cmd/run"
	"os"
)

var (
	rootConfigPath string
	rootConfig     *config.RootConfig = config.NewRootConfig()
	rootCommand                       = &cobra.Command{
		Use:     "power",
		Version: "1.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configErr := rootConfig.TryReadConfig(rootConfigPath)
			if configErr != nil {
				return configErr
			}

			return rootConfig.ResolveArgsAndFlags(cmd.Flags(), args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			value, _ := json.MarshalIndent(rootConfig, "", "  ")
			fmt.Printf("Executed root command. Resolved config: %s\n", value)
		},
	}
)

func init() {
	rootConfig.AttachFlagsToCommand(rootCommand)

	rootCommand.AddCommand(run.GetCommand())
}

// Execute is an entry-point function to start the CLI interactions
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
