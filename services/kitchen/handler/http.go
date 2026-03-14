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
		// Return newly created resource info
		writeError := util.WriteResponse(writer, http.StatusOK, response.Data)
		if writeError != nil {
			util.WriteError(writer, http.StatusInternalServerError, writeError)
		}
	}
}

func (oh *OrdersHttpHandler) ListOrders(writer http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()
	limit := util.GetQueryUint64(urlParams, "limit")
	if limit == nil {
		defaultLimit := uint64(10)
		limit = &defaultLimit
	}
	offset := util.GetQueryUint64(urlParams, "offset")
	if offset == nil {
		defaultOffset := uint64(0)
		offset = &defaultOffset
	}
	request := &orders.ListOrdersRequest{
		Limit:  limit,
		Offset: offset,
	}

	response, err := oh.service.ListOrders(request)
	if err != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	// Use consistent data/metadata structure
	metadata := &util.Metadata{
		Total: response.Meta.Total,
	}
	if limit != nil {
		metadata.Limit = limit
	}
	if offset != nil {
		metadata.Offset = offset
	}

	data := response.Data
	if data == nil {
		data = []*orders.Order{}
	}

	wrappedResponse := &util.ResponseWrapper{
		Data:     data,
		Metadata: metadata,
	}

	writeError := util.WriteResponse(writer, http.StatusOK, wrappedResponse)
	if writeError != nil {
		util.WriteError(writer, http.StatusInternalServerError, writeError)
	}
}
