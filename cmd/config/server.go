package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-grpc/internal/server"
)

type ServerConfig struct {
	*server.ServerConfig
	viperInstance *viper.Viper
	prefix        string
}

func NewServerConfig(viperInstance *viper.Viper, host string, port int, prefix string) *ServerConfig {
	return &ServerConfig{
		viperInstance: viperInstance,
		prefix:        prefix,
		ServerConfig: &server.ServerConfig{
			Host: host,
			Port: port,
		},
	}
}

func (sc *ServerConfig) RegisterFlags(cmd *cobra.Command) {
	sc.RegisterFlagsForFlagSet(cmd.Flags())
}

func (sc *ServerConfig) RegisterFlagsForFlagSet(flags *pflag.FlagSet) {
	hostKey := "host"
	portKey := "port"
	if sc.prefix != "" {
		hostKey = fmt.Sprintf("%s-host", sc.prefix)
		portKey = fmt.Sprintf("%s-port", sc.prefix)
	}

	flags.StringVar(&sc.Host, hostKey, sc.Host, "set server host")
	flags.IntVar(&sc.Port, portKey, sc.Port, "set server port")
}

func (sc *ServerConfig) ResolveFlagsAndArgs(flags *pflag.FlagSet, _ []string) error {
	var resolver ParamReader = NewDualReader(sc.viperInstance, flags)

	hostKey := "host"
	portKey := "port"
	if sc.prefix != "" {
		hostKey = fmt.Sprintf("%s.host", sc.prefix)
		portKey = fmt.Sprintf("%s.port", sc.prefix)
	}

	sc.Host = resolver.SafeGetString(hostKey, sc.Host)
	sc.Port = resolver.SafeGetInt(portKey, sc.Port)

	return nil
}
