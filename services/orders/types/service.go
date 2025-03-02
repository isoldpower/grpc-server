package types

import (
	"context"
	"golang-grpc/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
}

type OrdersHandlerType struct {
	OrdersService OrderService
	orders.UnimplementedOrderServiceServer
}

type OrderHandler interface {
	CreateOrder(context.Context, *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
}
