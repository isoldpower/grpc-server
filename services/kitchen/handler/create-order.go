package handler

import (
	"context"
	"golang-grpc/internal/util"
	"golang-grpc/services/common/genproto/orders"
	"net/http"
	"time"
)

func CreateOrderHandler(writer http.ResponseWriter, _ *http.Request) {
	response, err := CreateOrder()
	if err != nil {
		util.WriteError(writer, http.StatusInternalServerError, err)
	}

	writeError := util.WriteResponse(writer, http.StatusOK, response)
	if writeError != nil {
		util.WriteError(writer, http.StatusInternalServerError, writeError)
	}
}

func CreateOrder() (*orders.CreateOrderResponse, error) {
	connection := NewGRPCClient("localhost:3081")
	client := orders.NewOrderServiceClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	response, err := client.CreateOrder(ctx, &orders.CreateOrderRequest{
		CustomerID: 24,
		ProductID:  1,
		Quantity:   1,
	})

	cancel()
	return response, err
}
