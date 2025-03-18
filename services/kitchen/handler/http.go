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

func (oh *OrdersHttpHandler) CreateOrder(writer http.ResponseWriter, req *http.Request) {
	var requestBody orders.CreateOrderRequest
	err := util.ParseBody(req, &requestBody)
	if err != nil {
		util.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	request := &orders.CreateOrderRequest{
		CustomerID: requestBody.CustomerID,
		ProductID:  requestBody.ProductID,
		Quantity:   requestBody.Quantity,
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

func (oh *OrdersHttpHandler) ListOrders(writer http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()
	var limit = util.GetQueryUint64(urlParams, "limit")
	var offset = util.GetQueryUint64(urlParams, "offset")
	request := &orders.ListOrdersRequest{
		Limit:  limit,
		Offset: offset,
	}

	response, err := oh.service.ListOrders(request)
	if err != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	writeError := util.WriteResponse(writer, http.StatusOK, response)
	if writeError != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
	}
}
