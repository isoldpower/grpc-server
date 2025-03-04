package service

import (
	"context"
	"golang-grpc/services/common/genproto/orders"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

// TODO: Replace with database
var ordersStorage = make([]*orders.Order, 0)

func (s *OrderService) CreateOrder(
	_ context.Context,
	order *orders.Order,
) error {
	ordersStorage = append(ordersStorage, order)
	return nil
}

func (s *OrderService) GetOrdersList(
	_ context.Context,
) ([]*orders.Order, error) {
	return ordersStorage, nil
}
