package server

import "net/http"

type ServerRunConfig struct {
	WithGracefulShutdown bool
}

type ServerConfig struct {
	Port int
	Host string
}

type ServerRoute struct {
	Pattern string
	Handler http.HandlerFunc
}

type Server interface {
	Run(ServerRunConfig) error
	Stop() error
	GetDoneChannel() <-chan bool
}
