package db

import (
	"context"
	"fmt"

	"github.com/nico-phil/grpc-microservices/payment/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type Adapter struct {
	db *gorm.DB
}

type Payment struct {
	gorm.Model
	CustomerId int64
	OrderId int64
	TotalPrice float32
	Status string
}

func NewAdapter(dataSourceUrl string)(*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl))
	if openErr != nil {
		return nil, fmt.Errorf("db connection errror: %v", openErr)
	}

	// if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("payment"))); err != nil {
	// 	return nil, fmt.Errorf("db otel plugin error: %v", err)
	// }

	err := db.AutoMigrate(&Payment{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}


func(a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	paymentModel := Payment {
		CustomerId: payment.CustomerId,
		OrderId: payment.OrderId,
		TotalPrice: payment.TotalPrice,
		Status: payment.Status,
	}
	res := a.db.WithContext(ctx).Create(&paymentModel)
	if res.Error == nil {
		payment.ID = int64(paymentModel.ID)
	}
	return res.Error
}

func(a Adapter) Get(ctx context.Context, id string) (domain.Payment, error) {
	
	var paymentEntity Payment
	res := a.db.WithContext(ctx).First(&paymentEntity, id)
	
	payment := domain.Payment {
		ID: int64(paymentEntity.ID),
		CustomerId: paymentEntity.CustomerId,
		OrderId: paymentEntity.OrderId,
		TotalPrice: paymentEntity.TotalPrice,
		Status: paymentEntity.Status,
		CreatedAt: paymentEntity.CreatedAt.Unix(),
	}

	return payment, res.Error

}