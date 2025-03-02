package main

import (
	"google.golang.org/grpc"
	"log"
)

func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to grpc server: %v", err)
	}

	return conn
}

func main() {

}
