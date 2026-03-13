package main

import (
	"golang-grpc/cmd/config"
	cmdKitchen "golang-grpc/cmd/kitchen"
	"golang-grpc/internal/log"
	svcKitchen "golang-grpc/services/kitchen"

	"github.com/spf13/pflag"
)

func main() {
	rootConfig := config.NewRootConfig()
	kitchenConfig := cmdKitchen.NewKitchenConfig(rootConfig)

	flags := pflag.NewFlagSet("kitchen", pflag.ExitOnError)
	kitchenConfig.RegisterFlagsForFlagSet(flags)

	if err := flags.Parse(nil); err != nil {
		panic(err)
	} else if err := kitchenConfig.TryResolveConfig(""); err != nil {
		log.Warnln("Failed to resolve config: %v", err)
	} else if err := kitchenConfig.ResolveFlagsAndArgs(flags, nil); err != nil {
		panic(err)
	}

	log.Infoln("Starting Kitchen Service...")
	service := svcKitchen.NewKitchenService(kitchenConfig.Store)
	service.ExecuteExternal()
}
