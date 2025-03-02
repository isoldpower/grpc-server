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

type NetworkType string

const (
	NetworkTypeTCP         NetworkType = "tcp"
	NetworkTypeTCP4        NetworkType = "tcp4"
	NetworkTypeTCP6        NetworkType = "tcp6"
	NetworkTypeUnix        NetworkType = "unix"
	NetworkTypeUnixNetwork NetworkType = "unixpacket"
)

type HttpServerConfig struct {
	Network NetworkType
	ServerConfig
}

type HTTPServer struct {
	basicConfig *HttpServerConfig
	router      *http.ServeMux
	listener    net.Listener
	server      *http.Server
	doneChannel chan bool
}

func createListener(host string, port int, network NetworkType) (net.Listener, error) {
	serveAddress := fmt.Sprintf("%s:%d", host, port)
	listener, listenErr := net.Listen(string(network), serveAddress)
	if listenErr != nil {
		return listener, listenErr
	}

	return listener, nil
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
		doneChannel <- true
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
	<-ctx.Done()

	// Shut the server down in 5 seconds
	fmt.Println("\nShutting down gracefully, press Ctrl+C again to force")
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
func (hs *HTTPServer) Run() error {
	address := fmt.Sprintf("%s:%d", hs.basicConfig.Host, hs.basicConfig.Port)
	hs.server = &http.Server{
		Addr:    address,
		Handler: hs.router,
	}

	fmt.Printf("🔥 Listening at http://%s\n", address)
	go hs.serveRouter()
	go hs.trackGracefulShutdown()

	if <-hs.doneChannel {
		fmt.Println("🟢 Graceful shutdown complete.")
	} else {
		fmt.Println("❌ Exited with problems.")
	}

	return nil
}

func (hs *HTTPServer) Stop() error {
	err := hs.server.Close()

	if err != nil {
		fmt.Println("Error occurred while closing the HTTP server connection")
		fmt.Println(err.Error())
	}

	return err
}
