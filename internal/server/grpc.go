package server

import (
	"context"
	"fmt"
	"golang-grpc/internal/color"
	"golang-grpc/internal/log"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type GrpcServerConfig struct {
	ServerConfig
}

type GRPCServer struct {
	basicConfig       *GrpcServerConfig
	server            *grpc.Server
	listener          net.Listener
	serviceRegistrars []func(*grpc.Server)
	doneChannel       chan bool
	servingChannel    chan bool
}

func (gs *GRPCServer) trackGracefulShutdown() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-gs.doneChannel:
		log.Infoln("Internal server shutdown signal received")
		return
	case <-ctx.Done():
		log.Processln("Shutting down %s server gracefully", color.Green("gRPC"))
		log.RaiseLog(func() {
			log.Logln("%s Press %s again to force", log.GetIcon(log.AttentionIcon), color.Red("Ctrl+C"))
		})
		break
	}

	err := gs.Stop()
	if err != nil {
		log.PrintError("Error shutting down gracefully", err)
	}
}

// NewGRPCServer safely creates new GRPCServer instance with predefined private fields.
// It defines basic router, done channel and listener
func NewGRPCServer(basicConfig *GrpcServerConfig) *GRPCServer {
	doneChannel := make(chan bool, 1)
	listener, err := createListener(basicConfig.Host, basicConfig.Port, NetworkTypeTCP)
	if err != nil {
		log.PrintError("Failed to create gRPC listener", err)
		doneChannel <- false
	}

	return &GRPCServer{
		listener:          listener,
		basicConfig:       basicConfig,
		doneChannel:       doneChannel,
		serviceRegistrars: make([]func(*grpc.Server), 0),
		servingChannel:    make(chan bool, 1),
	}
}

// GetDoneChannel returns the read-only boolean channel with "done" indicator.
// The indicator signals whether the server finished its work.
func (gs *GRPCServer) GetDoneChannel() <-chan bool {
	return gs.doneChannel
}

// GetServingChannel returns the read-only boolean channel with "serving" indicator.
// The indicator signals whether the server is serving and accepting connections.
func (gs *GRPCServer) GetServingChannel() <-chan bool {
	return gs.servingChannel
}

// AddServiceRegistrar adds service register function to the list of
// functions. Later, the list is used to iterate and register all services
// at GRPCServer.Run method
func (gs *GRPCServer) AddServiceRegistrar(
	handler func(*grpc.Server),
) {
	gs.serviceRegistrars = append(gs.serviceRegistrars, handler)
}

// Run bootstrap the configured server and tracks whether
// the server was shut by user
func (gs *GRPCServer) Run(config ServerRunConfig) error {
	address := fmt.Sprintf("%s:%d", gs.basicConfig.Host, gs.basicConfig.Port)
	gs.server = grpc.NewServer()
	for _, registrar := range gs.serviceRegistrars {
		registrar(gs.server)
	}

	go func() {
		gs.servingChannel <- true
		serveError := gs.server.Serve(gs.listener)
		if serveError != nil {
			gs.doneChannel <- false
		}
	}()

	if !config.Silent {
		log.Processln("Listening at tcp://%s\n", address)
	}

	if config.WithGracefulShutdown {
		go gs.trackGracefulShutdown()
	}

	if <-gs.doneChannel {
		log.Successln("Graceful shutdown complete %s.", color.Green("(gRPC)"))
		gs.doneChannel <- true
	} else {
		log.Errorln("Exited with problems.")
		gs.doneChannel <- false
	}

	return nil
}

// Stop closes the server listener and sends the signal to the
// done channel that the server is closed
func (gs *GRPCServer) Stop() error {
	stopChannel := make(chan bool)

	go func() {
		gs.server.GracefulStop()
		close(stopChannel)
	}()

	select {
	case <-stopChannel:
		gs.doneChannel <- true
	case <-time.After(3 * time.Second):
		log.Infoln("Graceful shutdown timed out, forcing server shutdown")
		gs.server.Stop()
		gs.doneChannel <- false
	}

	return nil
}
