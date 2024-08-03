package domain

import "time"

type Payment struct {
	ID int64 `json:"id"`
	CustomerId  int64 `json:"customer_id"`
	OrderId int64 `json:"order_id"`
	TotalPrice float32 `json:"total_price"`
	Status string `json:"status"`
	CreatedAt int64 `json:"created_at"`
}

func NewPayment(curtomerId, orderId int64, totalPrice float32) Payment {
	return Payment {
		CustomerId: curtomerId,
		OrderId: orderId,
		TotalPrice: totalPrice,
		Status: "pending",
		CreatedAt: time.Now().Unix(), 
	}
}