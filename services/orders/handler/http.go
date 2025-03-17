package handler

import (
	"context"
	"golang-grpc/internal/server"
	"golang-grpc/internal/util"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"net/http"
)

type OrdersHttpHandler struct {
	types.OrdersHandlerType
}

func (oh *OrdersHttpHandler) tryCreateOrder(
	request *orders.CreateOrderRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *orders.Order) {
	order := &orders.Order{
		CustomerID: request.CustomerID,
		ProductID:  request.ProductID,
		Quantity:   request.Quantity,
	}

	createErr := oh.OrdersService.CreateOrder(context, order)
	if createErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, createErr)
		return false, order
	}

	return true, order
}

// NewHttpOrdersHandler creates an instance of OrdersHttpHandler object
// with specific predefined values to ensure safe usage and runtime
func NewHttpOrdersHandler(orderService types.OrderService) *OrdersHttpHandler {
	httpHandler := &OrdersHttpHandler{
		OrdersHandlerType: types.OrdersHandlerType{
			OrdersService: orderService,
		},
	}

	return httpHandler
}

// GetRoutes returns a list of HTTP routes related to the handler.
func (oh *OrdersHttpHandler) GetRoutes() []*server.ServerRoute {
	return []*server.ServerRoute{
		{Pattern: "POST /orders", Handler: oh.CreateOrder},
	}
}

// CreateOrder writes new order to local storage
func (oh *OrdersHttpHandler) CreateOrder(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var requestDto orders.CreateOrderRequest
	bodyErr := util.ParseBody(request, &requestDto)
	if bodyErr != nil {
		util.WriteError(writer, http.StatusBadRequest, bodyErr)
		return
	}

	if created, order := oh.tryCreateOrder(&requestDto, writer, request.Context()); created {
		resultErr := util.WriteResponse(writer, http.StatusCreated, order)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}
