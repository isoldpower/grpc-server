package main

import (
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/kitchen/handler"
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
	handler     *handler.OrdersHttpHandler
	server      *server.HTTPServer
}

func NewHTTPServer(basicConfig *httpServerConfig) *HTTPServer {
	return &HTTPServer{
		basicConfig: basicConfig,
		handler:     handler.NewOrdersHttpHandler(),
		server: server.NewHTTPServer(&server.HttpServerConfig{
			Network:      "tcp",
			ServerConfig: basicConfig.ServerConfig,
		}),
	}
}

func (hs *HTTPServer) registerRoutes() {
	hs.server.AddRoute("POST /", hs.handler.CreateOrder)
	hs.server.AddRoute("GET /", hs.handler.GetOrders)
}

func (hs *HTTPServer) Run(config server.ServerRunConfig) error {
	hs.registerRoutes()
	err := hs.server.Run(config)

	go func() {
		doneChannel := hs.server.GetDoneChannel()
		if <-doneChannel {
			fmt.Println("HTTP server shut down...")
		} else {
			fmt.Println("HTTP server forced to shut.")
		}
	}()

	return err
}

func (hs *HTTPServer) GetDoneChannel() <-chan bool {
	return hs.server.GetDoneChannel()
}

func (hs *HTTPServer) Stop() error {
	return hs.server.Stop()
}
