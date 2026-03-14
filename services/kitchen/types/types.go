package types

import "golang-grpc/services/common/genproto/orders"

type OrderService interface {
	CreateOrder(*orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
	ListOrders(request *orders.ListOrdersRequest) (*orders.ListOrdersResponse, error)
}
