package handler

import (
	"context"
	"errors"
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
	if request.CustomerID == "" {
		util.WriteError(writer, http.StatusBadRequest, errors.New("customer_id is required"))
		return false, nil
	}
	if request.ProductID == "" {
		util.WriteError(writer, http.StatusBadRequest, errors.New("product_id is required"))
		return false, nil
	}

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

func (oh *OrdersHttpHandler) tryCreateCustomer(
	request *orders.CreateCustomerRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *orders.Customer) {
	if request.Name == "" {
		util.WriteError(writer, http.StatusBadRequest, errors.New("name is required"))
		return false, nil
	}

	customer := &orders.Customer{
		Name: request.Name,
	}

	createErr := oh.OrdersService.CreateCustomer(context, customer)
	if createErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, createErr)
		return false, customer
	}

	return true, customer
}

func (oh *OrdersHttpHandler) tryCreateProduct(
	request *orders.CreateProductRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *orders.Product) {
	if request.Title == "" {
		util.WriteError(writer, http.StatusBadRequest, errors.New("title is required"))
		return false, nil
	}
	if request.Description == "" {
		util.WriteError(writer, http.StatusBadRequest, errors.New("description is required"))
		return false, nil
	}

	product := &orders.Product{
		Title:       request.Title,
		Description: request.Description,
	}

	createErr := oh.OrdersService.CreateProduct(context, product)
	if createErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, createErr)
		return false, product
	}

	return true, product
}

func (oh *OrdersHttpHandler) tryListOrders(
	req *orders.ListOrdersRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *util.ResponseWrapper) {
	listed, total, listErr := oh.OrdersService.GetOrdersList(req.Limit, req.Offset, context)
	if listErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, listErr)
		return false, nil
	}

	return true, &util.ResponseWrapper{
		Data: listed,
		Metadata: &util.Metadata{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
}

func (oh *OrdersHttpHandler) tryListCustomers(
	req *orders.ListCustomersRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *util.ResponseWrapper) {
	listed, total, listErr := oh.OrdersService.GetCustomersList(req.Limit, req.Offset, context)
	if listErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, listErr)
		return false, nil
	}

	return true, &util.ResponseWrapper{
		Data: listed,
		Metadata: &util.Metadata{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
}

func (oh *OrdersHttpHandler) tryListProducts(
	req *orders.ListProductsRequest,
	writer http.ResponseWriter,
	context context.Context,
) (bool, *util.ResponseWrapper) {
	listed, total, listErr := oh.OrdersService.GetProductsList(req.Limit, req.Offset, context)
	if listErr != nil {
		util.WriteError(writer, http.StatusInternalServerError, listErr)
		return false, nil
	}

	return true, &util.ResponseWrapper{
		Data: listed,
		Metadata: &util.Metadata{
			Total:  total,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
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
		{Pattern: "POST /customers", Handler: oh.CreateCustomer},
		{Pattern: "GET /customers", Handler: oh.GetCustomersList},
		{Pattern: "POST /products", Handler: oh.CreateProduct},
		{Pattern: "GET /products", Handler: oh.GetProductsList},
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

// CreateCustomer writes new customer to local storage
func (oh *OrdersHttpHandler) CreateCustomer(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var requestDto orders.CreateCustomerRequest
	bodyErr := util.ParseBody(request, &requestDto)
	if bodyErr != nil {
		util.WriteError(writer, http.StatusBadRequest, bodyErr)
		return
	}

	if created, customer := oh.tryCreateCustomer(&requestDto, writer, request.Context()); created {
		resultErr := util.WriteResponse(writer, http.StatusCreated, customer)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}

// CreateProduct writes new product to local storage
func (oh *OrdersHttpHandler) CreateProduct(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var requestDto orders.CreateProductRequest
	bodyErr := util.ParseBody(request, &requestDto)
	if bodyErr != nil {
		util.WriteError(writer, http.StatusBadRequest, bodyErr)
		return
	}

	if created, product := oh.tryCreateProduct(&requestDto, writer, request.Context()); created {
		resultErr := util.WriteResponse(writer, http.StatusCreated, product)
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
	limitVal, err := strconv.ParseUint(urlParams.Get("limit"), 10, 64)
	if err != nil || limitVal == 0 {
		limitVal = 10
	}
	limit := &limitVal

	offsetVal, err := strconv.ParseUint(urlParams.Get("offset"), 10, 64)
	if err != nil {
		offsetVal = 0
	}
	offset := &offsetVal

	requestDto := orders.ListOrdersRequest{
		Filters: &filters,
		Limit:   limit,
		Offset:  offset,
	}

	if success, response := oh.tryListOrders(&requestDto, writer, request.Context()); success {
		resultErr := util.WriteResponse(writer, http.StatusOK, response)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}

// GetCustomersList writes gets list of customers
func (oh *OrdersHttpHandler) GetCustomersList(
	writer http.ResponseWriter,
	request *http.Request,
) {
	urlParams := request.URL.Query()
	limitVal, err := strconv.ParseUint(urlParams.Get("limit"), 10, 64)
	if err != nil || limitVal == 0 {
		limitVal = 10
	}
	limit := &limitVal

	offsetVal, err := strconv.ParseUint(urlParams.Get("offset"), 10, 64)
	if err != nil {
		offsetVal = 0
	}
	offset := &offsetVal

	requestDto := orders.ListCustomersRequest{
		Limit:  limit,
		Offset: offset,
	}

	if success, response := oh.tryListCustomers(&requestDto, writer, request.Context()); success {
		resultErr := util.WriteResponse(writer, http.StatusOK, response)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}

// GetProductsList writes gets list of products
func (oh *OrdersHttpHandler) GetProductsList(
	writer http.ResponseWriter,
	request *http.Request,
) {
	urlParams := request.URL.Query()
	limitVal, err := strconv.ParseUint(urlParams.Get("limit"), 10, 64)
	if err != nil || limitVal == 0 {
		limitVal = 10
	}
	limit := &limitVal

	offsetVal, err := strconv.ParseUint(urlParams.Get("offset"), 10, 64)
	if err != nil {
		offsetVal = 0
	}
	offset := &offsetVal

	requestDto := orders.ListProductsRequest{
		Limit:  limit,
		Offset: offset,
	}

	if success, response := oh.tryListProducts(&requestDto, writer, request.Context()); success {
		resultErr := util.WriteResponse(writer, http.StatusOK, response)
		if resultErr != nil {
			util.WriteError(writer, http.StatusInternalServerError, resultErr)
		}
	}
}
