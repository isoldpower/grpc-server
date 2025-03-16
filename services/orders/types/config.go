package types

import (
	"golang-grpc/cmd/config"
)

type InitialConfig struct {
	Root *config.RootConfig
	Test string
}
