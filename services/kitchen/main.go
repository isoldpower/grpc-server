package main

import (
	"golang-grpc/internal/server"
	"google.golang.org/grpc"
	"log"
)

func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %v", err)
	}

	return conn
}

func main() {
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: server.ServerConfig{
			Port: 8000,
			Host: "localhost",
		},
	})

	httpServer.Run(server.ServerRunConfig{
		ReturnOnError: false,
	})
}
