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
		CustomerID: "0dacf448-068e-4182-ab9d-89537b43cb2e",
		ProductID:  "a568a8da-613f-4bc0-b10b-e645656562fa",
		Quantity:   2,
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
		CustomerID: "0dacf448-068e-4182-ab9d-89537b43cb2e",
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
