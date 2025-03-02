package main

import (
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/handler"
	"golang-grpc/services/orders/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GRPCServer struct {
	basicConfig *gRPCServerConfig
	server      *grpc.Server
}

type gRPCServerConfig struct {
	server.ServerConfig
}

func NewGRPCServer(defaultConfig *gRPCServerConfig) *GRPCServer {
	return &GRPCServer{basicConfig: defaultConfig}
}

func (s *GRPCServer) registerServices() {
	orderService := service.NewOrderService()
	handler.NewGrpcOrdersHandler(s.server, orderService)
}

func (s *GRPCServer) listenServer(errorChannel chan<- error) {
	listenAddress := fmt.Sprintf("%s:%d", s.basicConfig.Host, s.basicConfig.Port)
	tcpListener, tcpError := net.Listen(
		"tcp",
		listenAddress,
	)
	if tcpError != nil {
		errorChannel <- tcpError
		fmt.Printf("failed to listen: %v\n", tcpError)
	} else {
		fmt.Printf("Started gRPC server on tcp://%s\n", listenAddress)
	}

	s.server = grpc.NewServer()
	s.registerServices()

	serveError := s.server.Serve(tcpListener)
	if serveError != nil {
		errorChannel <- serveError
	}
}

func (s *GRPCServer) Run(_ server.ServerRunConfig) <-chan error {
	errorChannel := make(chan error)

	go s.listenServer(errorChannel)

	return errorChannel
}

func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()
	return nil
}
