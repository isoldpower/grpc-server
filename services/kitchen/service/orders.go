package service

import (
	"context"
	"fmt"
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
		fmt.Printf("failed to connect to grpc server: %v\n", err)
	}

	return grpcConnection
}

func NewOrderService() *OrderService {
	connection := createGrpcConnection("localhost:3081")

	return &OrderService{
		client: orders.NewOrderServiceClient(connection),
	}
}

func (og *OrderService) GetOrdersList(request *orders.GetOrdersRequest) (*orders.GetOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	response, err := og.client.GetOrdersList(ctx, request)

	cancel()
	return response, err
}

func (og *OrderService) CreateOrder(request *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	response, err := og.client.CreateOrder(ctx, request)

	cancel()
	return response, err
}
