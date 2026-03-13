package orders

import (
	"golang-grpc/internal/server"
	"golang-grpc/services/orders/storage"
	"golang-grpc/services/orders/types"
)

type OrdersService struct {
	config *types.InitialConfig
}

func NewOrdersService(config *types.InitialConfig) *OrdersService {
	return &OrdersService{
		config: config,
	}
}

func (os *OrdersService) ExecuteExternal() {
	ready := make(chan bool, 1)
	done := os.Execute(ready)
	<-done
}

func (os *OrdersService) Execute(ready chan<- bool) <-chan bool {
	storage.Config = os.config
	var grpcServer server.Server = NewGRPCServer(&gRPCServerConfig{
		ServerConfig: *os.config.GRPC,
	})
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: *os.config.HTTP,
	})

	runList := []server.Server{grpcServer, httpServer}
	wg := server.RunServersInParallel(runList, server.ServerRunConfig{
		WithGracefulShutdown: true,
		Silent:               true,
	})

	ready <- true
	doneChannel := make(chan bool, 1)
	go func() {
		wg.Wait()
		doneChannel <- true
	}()

	return doneChannel
}
