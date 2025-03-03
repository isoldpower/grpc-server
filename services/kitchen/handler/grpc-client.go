package handler

import (
	"fmt"
	"google.golang.org/grpc"
)

func NewGRPCClient(address string) *grpc.ClientConn {
	grpcConnection, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to connect to grpc server: %v\n", err)
	}

	return grpcConnection
}
