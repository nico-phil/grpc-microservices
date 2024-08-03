package domain

import (
	"time"
)

type OrderItem struct {
	ProductCode string `json:"product_code"`
	UnitPrice float32 `json:"unit_price"`
	Quantity int32 `json:"quantity"`
}

type Order struct {
	ID int64 `json:"id"`
	CustomerID int64 `json:"customer_id"`
	Status string `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt int64 `json:"created_at"`
}


func NewOder(customerId int64, orderItems []OrderItem) Order{
	return Order {
		CustomerID: customerId,
		Status: "Pending",
		OrderItems: orderItems,
		CreatedAt: time.Now().Unix(),
	}
}


func(o *Order) TotalPrice() float32{
	var totalPrice float32
	for _, orderItem := range o.OrderItems {
		 totalPrice += float32(orderItem.Quantity) * orderItem.UnitPrice 
	}

	return totalPrice
}
