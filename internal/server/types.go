package server

type ServerRunConfig struct {
	ReturnOnError bool
}

type ServerConfig struct {
	Port int
	Host string
}

type Server interface {
	Run(ServerRunConfig) <-chan error
	Stop() error
}
