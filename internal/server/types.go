package server

type ServerRunConfig struct {
}

type ServerConfig struct {
	Port int
	Host string
}

type Server interface {
	Run(ServerRunConfig) error
	Stop() error
}
