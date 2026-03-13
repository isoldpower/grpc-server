package kitchen

import (
	"golang-grpc/internal/server"
	"golang-grpc/services/kitchen/store"
)

type KitchenService struct {
	config *store.InitialConfig
}

func NewKitchenService(config *store.InitialConfig) *KitchenService {
	return &KitchenService{
		config: config,
	}
}

func (ks *KitchenService) ExecuteExternal() {
	ks.config = &store.InitialConfig{}

	ready := make(chan bool, 1)
	done := ks.Execute(ready)
	<-done
}

func (ks *KitchenService) Execute(ready chan<- bool) <-chan bool {
	var httpServer server.Server = NewHTTPServer(&httpServerConfig{
		ServerConfig: *ks.config.Server,
	})

	runList := []server.Server{httpServer}
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
