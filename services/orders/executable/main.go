package main

import (
	"golang-grpc/cmd/config"
	cmdOrders "golang-grpc/cmd/orders"
	"golang-grpc/internal/log"
	svcOrders "golang-grpc/services/orders"

	"github.com/spf13/pflag"
)

func main() {
	rootConfig := config.NewRootConfig()
	ordersConfig := cmdOrders.NewOrdersConfig(rootConfig)

	flags := pflag.NewFlagSet("orders", pflag.ExitOnError)
	ordersConfig.RegisterFlagsForFlagSet(flags)

	if err := flags.Parse(nil); err != nil {
		panic(err)
	}

	if err := ordersConfig.TryResolveConfig(""); err != nil {
		log.Warnln("Failed to resolve config: %v", err)
	}

	if err := ordersConfig.ResolveFlagsAndArgs(flags, nil); err != nil {
		panic(err)
	}

	log.Infoln("Starting Orders Service...")
	service := svcOrders.NewOrdersService(ordersConfig.Store)
	service.ExecuteExternal()
}
