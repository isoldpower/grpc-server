package handler

import (
	"golang-grpc/internal/util"
	"golang-grpc/services/common/genproto/orders"
	"golang-grpc/services/kitchen/service"
	"golang-grpc/services/kitchen/types"
	"net/http"
)

type OrdersHttpHandler struct {
	service types.OrderService
}

func NewOrdersHttpHandler() *OrdersHttpHandler {
	return &OrdersHttpHandler{
		service: service.NewOrderService(),
	}
}

func (oh *OrdersHttpHandler) CreateOrder(writer http.ResponseWriter, _ *http.Request) {
	request := &orders.CreateOrderRequest{
		CustomerID: 32,
		ProductID:  1,
		Quantity:   10,
	}

	if response, err := oh.service.CreateOrder(request); err != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
	} else {
		writeError := util.WriteResponse(writer, http.StatusOK, response)
		if writeError != nil {
			util.WriteError(writer, http.StatusInternalServerError, writeError)
		}
	}
}

func (oh *OrdersHttpHandler) GetOrders(writer http.ResponseWriter, _ *http.Request) {
	request := &orders.GetOrdersRequest{
		CustomerID: 32,
	}

	response, err := oh.service.GetOrdersList(request)
	if err != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	writeError := util.WriteResponse(writer, http.StatusOK, response)
	if writeError != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
	}
}
