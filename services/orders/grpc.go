package main

import (
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/handler"
	"golang-grpc/services/orders/service"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	basicConfig *gRPCServerConfig
	server      *server.GRPCServer
}

type gRPCServerConfig struct {
	server.ServerConfig
}

func (s *GRPCServer) registerServices() {
	s.server.AddServiceRegistrar(func(passed *grpc.Server) {
		handler.NewGrpcOrdersHandler(passed, service.NewOrderService())
	})
}

// NewGRPCServer creates new GRPCServer instance with basic settings applied.
// By default, it applies the list of handlers, creates server and saves config
func NewGRPCServer(defaultConfig *gRPCServerConfig) *GRPCServer {
	return &GRPCServer{
		basicConfig: defaultConfig,
		server: server.NewGRPCServer(&server.GrpcServerConfig{
			ServerConfig: defaultConfig.ServerConfig,
		}),
	}
}

// GetDoneChannel returns the boolean read-only channel with done signal.
// The transferred signal is true when the server shut down successfully and false when with errors
func (s *GRPCServer) GetDoneChannel() <-chan bool {
	return s.GetDoneChannel()
}

// Run starts the server to listen and handle at specific port.
// Returns possible server run process error
func (s *GRPCServer) Run(config server.ServerRunConfig) error {
	fmt.Println("🔄 Running gRPC server...")
	s.registerServices()

	return s.server.Run(config)
}

// Stop shuts down the server gracefully
func (s *GRPCServer) Stop() error {
	return s.server.Stop()
}
