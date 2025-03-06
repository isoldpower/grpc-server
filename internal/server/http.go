package server

import (
	"context"
	"errors"
	"fmt"
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
	basicConfig *HttpServerConfig
	router      *http.ServeMux
	listener    net.Listener
	server      *http.Server
	doneChannel chan bool
}

func (hs *HTTPServer) listenForErrors(errorChannel <-chan error) {
	for err := range errorChannel {
		fmt.Println("Error occurred while listening for errors: ")
		fmt.Println("\t", err)
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
		fmt.Printf("Failed to create listener: %v\n", err)
		doneChannel <- false
	}

	return &HTTPServer{
		router:      http.NewServeMux(),
		basicConfig: basicConfig,
		listener:    listener,
		doneChannel: doneChannel,
	}
}

// GetDoneChannel returns the read-only boolean channel with "done" indicator.
// The indicator signals whether the server finished its work.
func (hs *HTTPServer) GetDoneChannel() <-chan bool {
	return hs.doneChannel
}

func (hs *HTTPServer) trackGracefulShutdown() {
	// Track for shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-hs.doneChannel:
		fmt.Println("Internal server shutdown signal received")
		break
	case <-ctx.Done():
		fmt.Println("Shutting down HTTP server gracefully")
		fmt.Println("\t ↳ Press Ctrl+C again to force")
		break
	}

	// Shut the server down in 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Handle force shutdown
	if err := hs.server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown with error: %v", err)
		hs.doneChannel <- false
	}

	// Finish shutdown
	hs.doneChannel <- true
}

func (hs *HTTPServer) serveRouter() {
	err := hs.server.Serve(hs.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		hs.doneChannel <- false
		fmt.Println("Error occurred while serving: ", err)
	}
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

	fmt.Printf("🔥 Listening at http://%s\n", address)
	go hs.serveRouter()
	if config.WithGracefulShutdown {
		go hs.trackGracefulShutdown()
	}

	if <-hs.doneChannel {
		fmt.Println("🟢 Graceful shutdown complete.")
		hs.doneChannel <- true
	} else {
		fmt.Println("❌ Exited with problems.")
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
