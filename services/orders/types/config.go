package types

import (
	"golang-grpc/cmd/config"
	"golang-grpc/internal/database"
)

type InitialConfig struct {
	Root     *config.RootConfig
	Database *database.Config
}
