package types

import (
	"context"
	"golang-grpc/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrdersList(limit *uint64, offset *uint64, ctx context.Context) ([]*orders.Order, error)
}

type OrdersHandlerType struct {
	OrdersService OrderService
	orders.UnimplementedOrderServiceServer
}
