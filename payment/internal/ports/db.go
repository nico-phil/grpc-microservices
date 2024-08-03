package ports

import (
	"context"

	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/domain"
)

type DBPort interface {
	Save(context.Context,  domain.Payment) error
	Get(context.Context, string) (domain.Payment, error)
}