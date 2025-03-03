package handler

import (
	"context"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	types.OrdersHandlerType
}

// NewGrpcOrdersHandler registers new grpc handler to specified grpc server
func NewGrpcOrdersHandler(server *grpc.Server, orderService types.OrderService) {
	gRPCHandler := &OrdersGrpcHandler{
		OrdersHandlerType: types.OrdersHandlerType{
			OrdersService: orderService,
		},
	}

	orders.RegisterOrderServiceServer(server, gRPCHandler)
}

// CreateOrder writes new order to local storage from grpc request
func (h *OrdersGrpcHandler) CreateOrder(
	context context.Context,
	_ *orders.CreateOrderRequest,
) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		OrderID:    52,
		CustomerID: 2,
		ProductID:  1,
		Quantity:   1,
	}

	createError := h.OrdersService.CreateOrder(context, order)
	if createError != nil {
		return nil, createError
	}

	response := &orders.CreateOrderResponse{
		Status: "success",
	}

	return response, nil
}
