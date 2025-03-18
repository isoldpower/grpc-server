package handler

import (
	"context"
	"golang-grpc/internal/log"
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
	req *orders.CreateOrderRequest,
) (*orders.CreateOrderResponse, error) {
	order := &orders.Order{
		CustomerID: req.CustomerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	createError := h.OrdersService.CreateOrder(context, order)
	if createError != nil {
		log.PrintError("create order error occurred at gRPC connection", createError)
		return nil, createError
	}

	response := &orders.CreateOrderResponse{
		Status: orders.CreateStatus_ORDER_CREATED,
	}

	return response, nil
}

func (h *OrdersGrpcHandler) ListOrders(
	context context.Context,
	req *orders.ListOrdersRequest,
) (*orders.ListOrdersResponse, error) {
	var total uint64 = 5

	retrievedOrders, retrieveError := h.OrdersService.GetOrdersList(req.Offset, req.Limit, context)
	if retrieveError != nil {
		return nil, retrieveError
	}

	meta := &orders.ListMeta{
		Total: total,
	}
	if req.Offset != nil {
		meta.Offset = *req.Offset
	}
	if req.Limit != nil {
		meta.Limit = *req.Limit
	}
	response := &orders.ListOrdersResponse{
		Data: retrievedOrders,
		Meta: meta,
	}

	return response, nil
}
