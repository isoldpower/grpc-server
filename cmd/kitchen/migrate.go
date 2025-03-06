package kitchen

import "github.com/spf13/cobra"

var (
	migrateCommand = &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations process for kitchen databases",
		Long:  `This command relates to Kitchen Microservice CLI. It runs the further CLI process for resolving all configuration files and running all sorts of flexible migrations on developer's choice'.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {

}
