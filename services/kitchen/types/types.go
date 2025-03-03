package types

import "golang-grpc/services/common/genproto/orders"

type OrderService interface {
	CreateOrder(*orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
	GetOrdersList(*orders.GetOrdersRequest) (*orders.GetOrdersResponse, error)
}
