package service

import (
	"context"
	"errors"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/storage"
	"golang-grpc/services/orders/types"
)

type OrderService struct {
	storage         types.ObjectStore[orders.Order]
	customerStorage types.ObjectStore[orders.Customer]
	productStorage  types.ObjectStore[orders.Product]
}

func NewOrderService() *OrderService {
	return &OrderService{
		storage:         storage.NewPostgresStorage(storage.Config.Database),
		customerStorage: storage.NewCustomerPostgresStorage(storage.Config.Database),
		productStorage:  storage.NewProductPostgresStorage(storage.Config.Database),
	}
}

func (s *OrderService) CreateOrder(
	_ context.Context,
	order *orders.Order,
) error {
	err, _ := s.storage.AddItem(order)

	return err
}

func (s *OrderService) CreateCustomer(
	_ context.Context,
	customer *orders.Customer,
) error {
	err, _ := s.customerStorage.AddItem(customer)

	return err
}

func (s *OrderService) GetCustomersList(
	limit *uint64,
	offset *uint64,
	_ context.Context,
) ([]*orders.Customer, uint64, error) {
	listed, success := s.customerStorage.ListItems(limit, offset)
	if !success {
		return []*orders.Customer{}, 0, errors.New("error listing customers")
	}

	total, err := s.customerStorage.CountItems()
	if err != nil {
		return []*orders.Customer{}, 0, err
	}

	response := make([]*orders.Customer, len(listed))
	copy(response, listed)

	return response, total, nil
}

func (s *OrderService) CreateProduct(
	_ context.Context,
	product *orders.Product,
) error {
	err, _ := s.productStorage.AddItem(product)

	return err
}

func (s *OrderService) GetProductsList(
	limit *uint64,
	offset *uint64,
	_ context.Context,
) ([]*orders.Product, uint64, error) {
	listed, success := s.productStorage.ListItems(limit, offset)
	if !success {
		return []*orders.Product{}, 0, errors.New("error listing products")
	}

	total, err := s.productStorage.CountItems()
	if err != nil {
		return []*orders.Product{}, 0, err
	}

	response := make([]*orders.Product, len(listed))
	copy(response, listed)

	return response, total, nil
}

func (s *OrderService) GetOrdersList(
	limit *uint64,
	offset *uint64,
	_ context.Context,
) ([]*orders.Order, uint64, error) {
	total, countErr := s.storage.CountItems()
	if countErr != nil {
		return []*orders.Order{}, 0, countErr
	}

	listed, success := s.storage.ListItems(limit, offset)
	if !success {
		return []*orders.Order{}, 0, errors.New("failed to list items")
	}

	response := make([]*orders.Order, len(listed))
	copy(response, listed)

	return response, total, nil
}
