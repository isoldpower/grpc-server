package main

import (
	"github.com/spf13/pflag"
	"golang-grpc/cmd/config"
	"golang-grpc/internal/log"
	"golang-grpc/internal/util"
	"golang-grpc/services/orders"
	"golang-grpc/services/orders/types"
)

func main() {
	var configPath = util.ResolvePath("config_runtime.yaml", "")

	rootConfig := config.NewRootConfigAt(configPath)
	rootConfig.Cli.Debug = false
	rootConfig.Cli.NoIcons = true

	if err := rootConfig.ResolveFlagsAndArgs(
		pflag.NewFlagSet("", pflag.ContinueOnError),
		[]string{},
	); err != nil {
		panic(err)
	}
	if err := rootConfig.TryResolveConfig(configPath); err != nil {
		panic(err)
	}

	finalConfig := &types.InitialConfig{
		Root: rootConfig,
		Test: "Internal start!",
	}
	service := orders.NewOrdersService(finalConfig)

	log.Processln("Internal start of an %s Orders service", log.GetIcon(log.BoxIcon))
	log.Debugln("Config: %s\n", log.GetObjectPattern(finalConfig))
	service.ExecuteExternal()
}
