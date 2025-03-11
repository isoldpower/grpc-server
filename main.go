package main

import (
	"golang-grpc/cmd"
	"os"
)

func main() {
	cliRoot := cmd.NewCommand()
	if err := cliRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
