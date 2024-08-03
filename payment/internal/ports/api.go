package ports

import (
	"context"

	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/domain"
)

type APIPort interface {
	Charge(context.Context, domain.Payment) (domain.Payment, error)
}
