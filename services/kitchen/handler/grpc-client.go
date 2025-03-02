package handler

import (
	"google.golang.org/grpc"
	"log"
)

func NewGRPCClient(address string) *grpc.ClientConn {
	grpcConnection, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %v", err)
	}

	return grpcConnection
}
