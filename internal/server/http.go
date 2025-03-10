package server

import (
	"context"
	"errors"
	"fmt"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type HttpNetworkType string

const (
	NetworkTypeTCP         HttpNetworkType = "tcp"
	NetworkTypeTCP4        HttpNetworkType = "tcp4"
	NetworkTypeTCP6        HttpNetworkType = "tcp6"
	NetworkTypeUnix        HttpNetworkType = "unix"
	NetworkTypeUnixNetwork HttpNetworkType = "unixpacket"
)

type HttpServerConfig struct {
	Network HttpNetworkType
	ServerConfig
}

type HTTPServer struct {
	basicConfig    *HttpServerConfig
	router         *http.ServeMux
	listener       net.Listener
	server         *http.Server
	doneChannel    chan bool
	servingChannel chan bool
}

func (hs *HTTPServer) listenForErrors(errorChannel <-chan error) {
	for err := range errorChannel {
		log.PrintError("Error occurred while listening for errors", err)
	}
}

func (hs *HTTPServer) trackGracefulShutdown() {
	// Track for shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-hs.doneChannel:
		log.Infoln("Internal server shutdown signal received")
		break
	case <-ctx.Done():
		log.Processln("Shutting down %s server gracefully", color.Blue("HTTP"))
		log.RaiseLog(func() {
			log.Logln("%s Press %s again to force", log.GetIcon(log.AttentionIcon), color.Red("Ctrl+C"))
		})
		break
	}

	// Shut the server down in 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Handle force shutdown
	if err := hs.server.Shutdown(ctx); err != nil {
		log.PrintError("Server forced to shutdown with error", err)
		hs.doneChannel <- false
	}

	// Finish shutdown
	hs.doneChannel <- true
}

func (hs *HTTPServer) serveRouter() {
	hs.servingChannel <- true
	err := hs.server.Serve(hs.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		hs.doneChannel <- false
		log.PrintError("Error occurred while serving HTTP listener", err)
	}
}

// NewHTTPServer safely creates new HTTPServer instance with predefined private fields.
// It defines basic router, done channel and listener
func NewHTTPServer(basicConfig *HttpServerConfig) *HTTPServer {
	if basicConfig.Network == "" {
		basicConfig.Network = NetworkTypeTCP
	}

	doneChannel := make(chan bool, 1)
	listener, err := createListener(basicConfig.Host, basicConfig.Port, basicConfig.Network)
	if err != nil {
		log.PrintError("Failed to create HTTP listener", err)
		doneChannel <- false
	}

	return &HTTPServer{
		router:         http.NewServeMux(),
		basicConfig:    basicConfig,
		listener:       listener,
		doneChannel:    doneChannel,
		servingChannel: make(chan bool, 1),
	}
}

// GetDoneChannel returns the read-only boolean channel with "done" indicator.
// The indicator signals whether the server finished its work.
func (hs *HTTPServer) GetDoneChannel() <-chan bool {
	return hs.doneChannel
}

// GetServingChannel returns the read-only boolean channel with "serving" indicator.
// The indicator signals whether the server is serving and accepting connections.
func (hs *HTTPServer) GetServingChannel() <-chan bool {
	return hs.servingChannel
}

// AddRoute adds additional route to server's http.ServeMux handler
func (hs *HTTPServer) AddRoute(
	pattern string,
	handler func(http.ResponseWriter, *http.Request),
) {
	hs.router.HandleFunc(pattern, handler)
}

// Run bootstrap the configured server and tracks whether
// the server was shut by user
func (hs *HTTPServer) Run(config ServerRunConfig) error {
	address := fmt.Sprintf("%s:%d", hs.basicConfig.Host, hs.basicConfig.Port)
	hs.server = &http.Server{
		Addr:    address,
		Handler: hs.router,
	}

	if !config.Silent {
		log.Processln("Listening at http://%s\n", address)
	}
	go hs.serveRouter()

	if config.WithGracefulShutdown {
		go hs.trackGracefulShutdown()
	}

	if <-hs.doneChannel {
		log.Successln("Graceful shutdown complete %s.", color.Blue("(HTTP)"))
		hs.doneChannel <- true
	} else {
		log.Errorln("Exited with problems.")
		hs.doneChannel <- false
	}

	return nil
}

// Stop closes the server listener and sends the signal to the
// done channel that the server is closed
func (hs *HTTPServer) Stop() error {
	hs.doneChannel <- true

	return nil
}
