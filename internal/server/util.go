package server

import (
	"fmt"
	"golang-grpc/internal/log"
	"net"
	"sync"
)

func createListener(host string, port int, network HttpNetworkType) (net.Listener, error) {
	serveAddress := fmt.Sprintf("%s:%d", host, port)
	listener, listenErr := net.Listen(string(network), serveAddress)
	if listenErr != nil {
		return listener, listenErr
	}

	return listener, nil
}

func RunServersInParallel(servers []Server, runConfig ServerRunConfig) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	for _, server := range servers {
		go func() {
			defer wg.Done()
			runError := server.Run(runConfig)

			if runError != nil {
				log.PrintError("Failed to serve the http server", runError)
				panic(runError)
			}
		}()

	}

	for _, server := range servers {
		<-server.GetServingChannel()
	}

	return &wg
}
