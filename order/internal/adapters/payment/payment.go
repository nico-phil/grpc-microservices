package payment

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/nico-phil/grpc-microservices/order/internal/application/core/domain"
	"github.com/nico-phil/microservices-proto/golang/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	paymentClient payment.PaymentClient
}

func NewAdapter(paymentServiceUrl  string) (*Adapter, error) {
	var opts [] grpc.DialOption

	opts = append(opts, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithCodes(
					codes.Unavailable, codes.ResourceExhausted,
		), grpc_retry.WithMax(5), 
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)))),
	)

	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	
	client := payment.NewPaymentClient(conn)
	
	return &Adapter{paymentClient: client}, nil
}

// why pass a pointer type order
func(a Adapter) Charge(ctx context.Context, order *domain.Order) error {
	 _, err := a.paymentClient.Create(ctx, &payment.CreatePaymentRequest{
		OrderId: order.ID,
		UserId: order.CustomerID,
		TotalPrice: order.TotalPrice(),
	})
	return err	
}

