package service

import (
	"context"
	"errors"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/storage"
	"golang-grpc/services/orders/types"
)

type OrderService struct {
	storage types.ObjectStore[orders.Order]
}

func NewOrderService() *OrderService {
	// TODO: Read config from global context
	return &OrderService{
		storage: storage.NewPostgresStorage(storage.Config.Database),
	}
}

func (s *OrderService) CreateOrder(
	_ context.Context,
	order *orders.Order,
) error {
	err, _ := s.storage.AddItem(order)

	return err
}

func (s *OrderService) GetOrdersList(
	_ context.Context,
) ([]*orders.Order, error) {
	listed, success := s.storage.ListItems(0, 10)
	if !success {
		return []*orders.Order{}, errors.New("failed to list items")
	}

	response := make([]*orders.Order, len(listed))
	for i, item := range listed {
		response[i] = item
	}

	return response, nil
}
