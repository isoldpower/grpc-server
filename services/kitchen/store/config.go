package store

import (
	"golang-grpc/cmd/config"
	"golang-grpc/internal/server"
)

type InitialConfig struct {
	Root   *config.RootConfig
	Server *server.ServerConfig
	Test   string
}
