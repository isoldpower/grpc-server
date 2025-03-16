package main

import (
	"github.com/spf13/pflag"
	"golang-grpc/cmd/config"
	"golang-grpc/internal/database"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
	"golang-grpc/services/orders"
	"golang-grpc/services/orders/types"
	"os"
	"strconv"
)

var (
	configPath  = util.ResolvePath("config_runtime.yaml", "")
	isDebug     = false
	isWithIcons = true
)

func createRootConfig() *types.InitialConfig {
	rootConfig := config.NewRootConfigAt(configPath)
	rootConfig.Cli.Debug = isDebug
	rootConfig.Cli.NoIcons = !isWithIcons

	if err := rootConfig.ResolveFlagsAndArgs(
		pflag.NewFlagSet("", pflag.ContinueOnError),
		[]string{},
	); err != nil {
		panic(err)
	}
	if err := rootConfig.TryResolveConfig(configPath); err != nil {
		panic(err)
	}

	return &types.InitialConfig{
		Root: rootConfig,
	}
}

func createDatabaseConfig() *database.Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort, portErr := strconv.Atoi(os.Getenv("DB_PORT"))
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbSchema := os.Getenv("DB_SCHEMA")

	if dbHost == "" {
		dbHost = "localhost"
		log.Warnln("DB_HOST environment variable not set. Falling back to %s", dbHost)
	}
	if portErr != nil || dbPort == 0 {
		dbPort = 5432
		log.Warnln("DB_PORT not set. Falling back to %d", dbPort)
	}
	if dbUsername == "" {
		dbUsername = "postgres"
		log.Warnln("DB_USERNAME environment variable not set. Falling back to %s", dbUsername)
	}
	if dbPassword == "" {
		dbPassword = "<empty string>"
		log.Warnln("DB_PASSWORD environment variable not set. Falling back to %s", dbPassword)
	}
	if dbDatabase == "" {
		dbDatabase = "orders"
		log.Warnln("DB_DATABASE environment variable not set. Falling back to %s", dbDatabase)
	}
	if dbSchema == "" {
		dbSchema = "PUBLIC"
		log.Warnln("DB_SCHEMA environment variable not set. Falling back to %s", dbSchema)
	}

	return &database.Config{
		Host:     dbHost,
		Port:     dbPort,
		Username: dbUsername,
		Password: dbPassword,
		Database: dbDatabase,
		Schema:   dbSchema,
	}
}

func main() {
	finalConfig := createRootConfig()
	finalConfig.Database = createDatabaseConfig()

	service := orders.NewOrdersService(finalConfig)

	log.Processln("Internal start of an %s Orders service", log.GetIcon(log.BoxIcon))
	log.Debugln("Config: %s\n", log.GetObjectPattern(finalConfig))
	service.ExecuteExternal()
}
