package handler

import (
	"context"
	"golang-grpc/internal/util"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"net/http"
)

type OrdersHttpHandler struct {
	types.OrdersHandlerType
}

func NewHttpOrdersHandler(orderService types.OrderService) *OrdersHttpHandler {
	httpHandler := &OrdersHttpHandler{
		OrdersHandlerType: types.OrdersHandlerType{
			OrdersService: orderService,
		},
	}

	return httpHandler
}

func (oh *OrdersHttpHandler) tryCreateOrder(
	request *orders.CreateOrderRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *orders.Order) {
	order := &orders.Order{
		OrderID:    42,
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

func (oh *OrdersHttpHandler) RegisterRouter(router *http.ServeMux) {
	router.HandleFunc("POST /orders", oh.CreateOrder)
}

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
