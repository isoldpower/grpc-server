package main

import (
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/handler"
	"golang-grpc/services/orders/service"
	"log"
	"net"
	"net/http"
)

type httpServerConfig struct {
	server.ServerConfig
}

type HTTPServer struct {
	basicConfig *httpServerConfig
	router      *http.ServeMux
	listener    net.Listener
}

func NewHTTPServer(basicConfig *httpServerConfig) *HTTPServer {
	return &HTTPServer{
		basicConfig: basicConfig,
	}
}

func (hs *HTTPServer) startServer(errorChannel chan<- error) {
	serveAddress := fmt.Sprintf("%s:%d", hs.basicConfig.Host, hs.basicConfig.Port)
	listener, _ := net.Listen("tcp", serveAddress)

	hs.listener = listener
	fmt.Printf("Started HTTP server on http://%s\n", serveAddress)

	serveError := http.Serve(listener, hs.router)
	errorChannel <- serveError
}

func (hs *HTTPServer) listenForErrors(errorChannel <-chan error) {
	for err := range errorChannel {
		fmt.Println("Error occurred while listening for errors: ")
		fmt.Println("\t", err)
	}
}

func (hs *HTTPServer) Run(_ server.ServerRunConfig) <-chan error {
	errorChannel := make(chan error)
	hs.router = http.NewServeMux()

	orderService := service.NewOrderService()
	orderHandler := handler.NewHttpOrdersHandler(orderService)
	orderHandler.RegisterRouter(hs.router)

	go hs.startServer(errorChannel)
	go hs.listenForErrors(errorChannel)

	return errorChannel
}

func (hs *HTTPServer) Stop() error {
	err := hs.listener.Close()
	if err != nil {
		fmt.Println("Error occurred while closing the HTTP connection")
		fmt.Println(err.Error())
	}

	return err
}
