package kitchen

import (
	"fmt"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
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
	log.Processln(
		"Running %s Kitchen HTTP server at %s",
		log.GetIcon(log.KnifeIcon),
		color.Green(fmt.Sprintf("%s:%d", hs.basicConfig.Host, hs.basicConfig.Port)),
	)
	log.RaiseLog(func() {
		log.Logln(
			"%s Press %s to exit",
			log.GetIcon(log.AttentionIcon),
			color.Red("Ctrl+C"),
		)
	})
	hs.registerRoutes()
	err := hs.server.Run(config)

	go func() {
		doneChannel := hs.server.GetDoneChannel()
		if <-doneChannel {
			log.Infoln("%s Kitchen HTTP server shut down...", log.GetIcon(log.KnifeIcon))
		} else {
			log.Errorln("%s Kitchen HTTP server forced to shut.", log.GetIcon(log.KnifeIcon))
		}
	}()

	return err
}

func (hs *HTTPServer) GetDoneChannel() <-chan bool {
	return hs.server.GetDoneChannel()
}

func (hs *HTTPServer) GetServingChannel() <-chan bool {
	return hs.server.GetServingChannel()
}

func (hs *HTTPServer) Stop() error {
	return hs.server.Stop()
}
