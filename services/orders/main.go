package orders

import (
	"golang-grpc/internal/server"
)

var (
	grpcServer server.Server = NewGRPCServer(&gRPCServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3081,
		},
	})
	httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3082,
		},
	})
)

func StartOrdersService() {
	runList := []server.Server{grpcServer, httpServer}
	waitServers := server.RunServersInParallel(runList, server.ServerRunConfig{
		WithGracefulShutdown: true,
	})

	waitServers.Wait()
}
