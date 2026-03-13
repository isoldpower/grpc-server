package types

import (
	"context"
	"golang-grpc/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrdersList(limit *uint64, offset *uint64, ctx context.Context) ([]*orders.Order, uint64, error)
	CreateCustomer(context.Context, *orders.Customer) error
	GetCustomersList(limit *uint64, offset *uint64, ctx context.Context) ([]*orders.Customer, uint64, error)
	CreateProduct(context.Context, *orders.Product) error
	GetProductsList(limit *uint64, offset *uint64, ctx context.Context) ([]*orders.Product, uint64, error)
}

type OrdersHandlerType struct {
	OrdersService OrderService
	orders.UnimplementedOrderServiceServer
}
