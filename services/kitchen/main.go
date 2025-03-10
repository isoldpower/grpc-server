package kitchen

import (
	"fmt"
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
		ServerConfig: server.ServerConfig{
			Port: 8000,
			Host: "localhost",
		},
	})

	go func() {
		err := httpServer.Run(server.ServerRunConfig{
			WithGracefulShutdown: true,
			Silent:               true,
		})

		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	ready <- <-httpServer.GetServingChannel()
	return httpServer.GetDoneChannel()
}
