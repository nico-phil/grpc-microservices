package db

import (
	"context"
	"fmt"

	"github.com/nico-phil/grpc-microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}


type Order struct {
	gorm.Model
	CustomerID int64
	Status string
	OrderItems [] OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice float32
	Quantity int32
	OrderId uint
}



func NewAdapter(dataSourceUrl string)(*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl))
	if openErr != nil {
		return nil, fmt.Errorf("db connection errror: %v", openErr)
	}

	err := db.AutoMigrate(&Order{}, OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Get(ctx context.Context, id int64) (domain.Order, error) {
	var orderEntity Order
	res := a.db.WithContext(ctx).Preload("OrderItems").First(&orderEntity, id)

	var orderItems []domain.OrderItem

	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice: orderItem.UnitPrice,
			Quantity: orderItem.Quantity,
		})
	}

	order := domain.Order {
		ID: int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status: orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt: orderEntity.CreatedAt.UnixNano(),
	}
	return order, res.Error
	
}


func (a Adapter) Save(ctx context.Context, order *domain.Order) error {
	var orderItems []OrderItem

	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice: orderItem.UnitPrice,
			Quantity: orderItem.Quantity,
		})
	}

	orderModel := Order {
		CustomerID: order.CustomerID,
		Status: order.Status,
		OrderItems: orderItems,
	}

	res := a.db.Create(&orderModel)
	if res.Error  == nil {
		order.ID = int64(orderModel.ID)
	}
	return res.Error
}