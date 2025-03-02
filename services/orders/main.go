package main

import (
	"golang-grpc/internal/server"
	"log"
)

func main() {
	var grpcServer server.Server = NewGRPCServer(&gRPCServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3081,
		},
	})
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: server.ServerConfig{
			Host: "localhost",
			Port: 3082,
		},
	})

	grpcError := grpcServer.Run(server.ServerRunConfig{
		ReturnOnError: true,
	})
	httpError := httpServer.Run(server.ServerRunConfig{
		ReturnOnError: true,
	})

	select {
	case err := <-grpcError:
		fmt.Println("Error occurred while listening at the gRPC connection")
		fmt.Println(err.Error())
		break
	case err := <-httpError:
		fmt.Println("Error occurred while listening at the HTTP connection")
		fmt.Println(err.Error())
		break
	}
}
