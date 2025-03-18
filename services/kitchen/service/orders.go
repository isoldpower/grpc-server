package service

import (
	"context"
	"fmt"
	"golang-grpc/internal/log"
	"golang-grpc/services/common/genproto/orders"
	"google.golang.org/grpc"
	"time"
)

type OrderService struct {
	client orders.OrderServiceClient
}

func createGrpcConnection(address string) *grpc.ClientConn {
	grpcConnection, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		log.PrintError(
			fmt.Sprintf("Failed to connect to grpc server at: %s", address),
			err,
		)
	}

	return grpcConnection
}

func NewOrderService() *OrderService {
	connection := createGrpcConnection("localhost:3081")

	return &OrderService{
		client: orders.NewOrderServiceClient(connection),
	}
}

func (og *OrderService) ListOrders(request *orders.ListOrdersRequest) (*orders.ListOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	response, err := og.client.ListOrders(ctx, request)

	cancel()
	return response, err
}

func (og *OrderService) CreateOrder(request *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	response, err := og.client.CreateOrder(ctx, request)

	cancel()
	return response, err
}
