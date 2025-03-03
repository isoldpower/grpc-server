package service

import (
	"context"
	ordersCommon "golang-grpc/services/common/genproto/orders"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

// TODO: Replace with database
var orders = make([]*ordersCommon.Order, 0)

func (s *OrderService) CreateOrder(
	_ context.Context,
	order *ordersCommon.Order,
) error {
	orders = append(orders, order)
	return nil
}

func (s *OrderService) GetOrdersList(
	_ context.Context,
) ([]*ordersCommon.Order, error) {
	return orders, nil
}
