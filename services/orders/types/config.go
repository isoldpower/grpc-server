package types

import (
	"golang-grpc/cmd/config"
	"golang-grpc/internal/database"
	"golang-grpc/internal/server"
)

type InitialConfig struct {
	Root     *config.RootConfig
	Database *database.Config
	GRPC     *server.ServerConfig
	HTTP     *server.ServerConfig
}
