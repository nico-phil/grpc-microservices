package grpc

import (
	"context"
	"fmt"

	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/domain"
	"github.com/nico-phil/microservices-proto/golang/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func(a Adapter) Create(ctx context.Context, request *payment.CreatePaymentRequest)(*payment.CreatePaymentResponse, error){
	newPayment := domain.NewPayment(request.UserId, request.OrderId, request.TotalPrice)

	result, err := a.api.Charge(ctx, newPayment)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge. %v", err)).Err()
	}

	fmt.Println("result",result.ID)

	return &payment.CreatePaymentResponse{BillId: 1}, nil
}