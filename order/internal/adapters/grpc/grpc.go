package grpc

import (
	"context"
	"log/slog"

	"github.com/nico-phil/grpc-microservices/order/internal/application/core/domain"
	"github.com/nico-phil/microservices-proto/golang/order"
)

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest)(*order.CreateOrderResponse, error) {
	slog.Info("creating order....")
	var orderItems []domain.OrderItem
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice: orderItem.UnitPrice,
			Quantity: orderItem.Quantity,
		})
	}

	newOrder := domain.NewOder(request.UserId, orderItems)
	result, err := a.api.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: result.ID}, nil
}

func (a Adapter) Get(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	result, err := a.api.GetOrder(ctx, request.OrderId)
	var orderItems []*order.OrderItem
	for _, orderItem := range result.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	if err != nil {
		return nil, err
	}
	return &order.GetOrderResponse{UserId: result.CustomerID, OrderItems: orderItems}, nil
}

func (a Adapter) GetTest(ctx context.Context, request *order.GetTestRequest) (*order.GetTestReponse, error) {
	return  &order.GetTestReponse{}, nil
}
