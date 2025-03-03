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

	defer httpServer.Stop()
	err := httpServer.Run(server.ServerRunConfig{
		WithGracefulShutdown: true,
	})
	if err != nil {
		fmt.Fatal(err)
	}
}
