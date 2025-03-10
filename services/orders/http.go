package orders

import (
	"fmt"
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/handler"
	"golang-grpc/services/orders/service"
	"net"
	"net/http"
)

type httpServerConfig struct {
	server.ServerConfig
}

type HTTPServer struct {
	basicConfig *httpServerConfig
	router      *http.ServeMux
	server      *server.HTTPServer
	handlers    []*handler.OrdersHttpHandler
	listener    net.Listener
}

func (hs *HTTPServer) registerRoutes() {
	for _, serverHandler := range hs.handlers {
		handlerRoutes := serverHandler.GetRoutes()
		for _, route := range handlerRoutes {
			hs.server.AddRoute(route.Pattern, route.Handler)
		}
	}
}

// NewHTTPServer creates new HTTPServer instance with basic settings applied.
// By default, it applies the list of handlers, creates server and saves config
func NewHTTPServer(basicConfig *httpServerConfig) *HTTPServer {
	return &HTTPServer{
		basicConfig: basicConfig,
		handlers: []*handler.OrdersHttpHandler{
			handler.NewHttpOrdersHandler(service.NewOrderService()),
		},
		server: server.NewHTTPServer(&server.HttpServerConfig{
			ServerConfig: basicConfig.ServerConfig,
		}),
	}
}

// Run bootstraps the orders HTTP server with desired logging
func (hs *HTTPServer) Run(config server.ServerRunConfig) error {
	fmt.Println("🔄 Running 📦Orders HTTP server...")
	hs.registerRoutes()
	return hs.server.Run(config)
}

// GetDoneChannel returns the channel with done signal.
// Signal is true if the server finished successfully and false if server finished with error
func (hs *HTTPServer) GetDoneChannel() <-chan bool {
	return hs.server.GetDoneChannel()
}

// GetServingChannel returns the read-only boolean channel with "serving" indicator.
// The indicator signals whether the server is serving and accepting connections.
func (hs *HTTPServer) GetServingChannel() <-chan bool {
	return hs.server.GetServingChannel()
}

// Stop gracefully closes the connection to the server
func (hs *HTTPServer) Stop() error {
	err := hs.listener.Close()
	if err != nil {
		fmt.Println("Error occurred while closing the HTTP connection")
		fmt.Println(err.Error())
	}

	return err
}
