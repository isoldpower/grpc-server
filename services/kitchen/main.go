package main

import (
	"golang-grpc/internal/server"
	"log"
)

func main() {
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: server.ServerConfig{
			Port: 8000,
			Host: "localhost",
		},
	})

	err := httpServer.Run(server.ServerRunConfig{})
	if err != nil {
		log.Fatal(err)
	}
}
