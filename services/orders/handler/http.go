package handler

import (
	"context"
	"golang-grpc/internal/server"
	"golang-grpc/internal/util"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/orders/types"
	"net/http"
	"strconv"
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

func (oh *OrdersHttpHandler) tryListOrders(
	req *orders.ListOrdersRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, []*orders.Order) {
	if listed, listErr := oh.OrdersService.GetOrdersList(req.Limit, req.Offset, context); listErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, listErr)
		return false, []*orders.Order{}
	} else {
		return true, listed
	}
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
		{Pattern: "GET /orders", Handler: oh.GetOrdersList},
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

// GetOrdersList writes gets list of orders
func (oh *OrdersHttpHandler) GetOrdersList(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var filters orders.Order
	if bodyErr := util.ParseBody(request, &filters); bodyErr != nil {
		util.WriteError(writer, http.StatusBadRequest, bodyErr)
		return
	}

	urlParams := request.URL.Query()
	var limit *uint64 = nil
	var offset *uint64 = nil
	if lim, err := strconv.ParseUint(urlParams.Get("limit"), 10, 64); err == nil && lim != 0 {
		limit = &lim
	}
	if off, err := strconv.ParseUint(urlParams.Get("offset"), 10, 64); err == nil && off != 0 {
		offset = &off
	}

	requestDto := orders.ListOrdersRequest{
		Filters: &filters,
		Limit:   limit,
		Offset:  offset,
	}

	if created, order := oh.tryListOrders(&requestDto, writer, request.Context()); created {
		resultErr := util.WriteResponse(writer, http.StatusOK, order)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}
