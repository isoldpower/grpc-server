package kitchen

import (
	"fmt"
	"golang-grpc/internal/server"
)

func StartKitchenService() {
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
		fmt.Printf("%v\n", err)
	}
}
