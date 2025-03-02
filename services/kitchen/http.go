package main

import (
	"context"
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/common/genproto/orders"
	"log"
	"net"
	"net/http"
	"time"
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

func (hs *HTTPServer) startServer() error {
	serveAddress := fmt.Sprintf("%s:%d", hs.basicConfig.Host, hs.basicConfig.Port)
	listener, listenErr := net.Listen("tcp", serveAddress)
	if listenErr != nil {
		return listenErr
	}

	hs.listener = listener
	log.Printf("Started HTTP server on http://%s\n", serveAddress)

	serveError := http.Serve(listener, hs.router)
	return serveError
}

func (hs *HTTPServer) listenForErrors(errorChannel <-chan error) {
	for err := range errorChannel {
		log.Println("Error occurred while listening for errors: ")
		log.Println("\t", err)
	}
}

func (hs *HTTPServer) Run(_ server.ServerRunConfig) <-chan error {
	hs.router = http.NewServeMux()

	interruptChannel := make(chan error)

	connection := NewGRPCClient("localhost:3081")
	hs.router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		client := orders.NewOrderServiceClient(connection)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		response, err := client.CreateOrder(ctx, &orders.CreateOrderRequest{
			CustomerID: 24,
			ProductID:  1,
			Quantity:   1,
		})

		if err != nil {
			cancel()
			log.Fatalf("client error: %v\n", err)
		}

		log.Printf("CreateOrderResponse: %v\n", response)
		cancel()
	})

	startError := hs.startServer()
	if startError != nil {
		interruptChannel <- startError
	}

	return interruptChannel
}

func (hs *HTTPServer) Stop() error {
	err := hs.listener.Close()
	if err != nil {
		log.Println("Error occurred while closing the HTTP connection")
		log.Println(err.Error())
	}

	return err
}
