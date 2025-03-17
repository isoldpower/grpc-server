package orders

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type MigrateCommand struct {
	config          *Config
	commandInstance *cobra.Command
}

func NewMigrateCommand(config *Config) *MigrateCommand {
	availableCommand := []string{"up", "status", "down", "up-by-one", "down-by-one"}

	return &MigrateCommand{
		config: config,
		commandInstance: &cobra.Command{
			Use:   fmt.Sprintf("migrate [%s]", strings.Join(availableCommand, " / ")),
			Short: "Migrate Orders microservice",
			Long:  "Run further CLI process of migrating all Orders services",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return util.ProtectedAction(cmd.Parent().PreRunE(cmd, args), func() error {
					return nil
				})
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				log.Infoln("Executed %s command", color.Underline("migrate orders"))

				command := "status"
				if len(args) > 0 && slices.Contains(availableCommand, args[0]) {
					command = args[0]
				}
				return run(config, command)
			},
		},
	}
}

func run(config *Config, command string) error {
	databaseLine := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.databaseConfig.Host,
		config.databaseConfig.Port,
		config.databaseConfig.Username,
		config.databaseConfig.Password,
		config.databaseConfig.Database,
		"disable",
		config.databaseConfig.Schema,
	)
	migrationDir := filepath.Join(util.RootPath(), "services", "orders", "_migrations")
	log.Debugln("database settings line: %s", databaseLine)
	log.Debugln("migration files directory: %s", migrationDir)

	gooseCommand := exec.Command("goose", "-dir", migrationDir, "postgres", databaseLine, command)
	if output, err := gooseCommand.CombinedOutput(); err != nil {
		log.PrintError("Error running goose on postgres", err)
		return nil
	} else {
		log.Infoln(string(output))
		return err
	}
}

func (mc *MigrateCommand) Register(parentCmd *cobra.Command) {
	parentCmd.AddCommand(mc.commandInstance)
}
