package models

import "time"

type OrderStatus string

const (
	OrderStatusAccepted  OrderStatus = "accepted"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusOnTheWay  OrderStatus = "on_the_way"
	OrderStatusDelivered OrderStatus = "delivered"
)

var orderStatusFlow = map[OrderStatus]OrderStatus{
	OrderStatusAccepted:  OrderStatusPreparing,
	OrderStatusPreparing: OrderStatusOnTheWay,
	OrderStatusOnTheWay:  OrderStatusDelivered,
}

func (s OrderStatus) Next() (OrderStatus, bool) {
	next, ok := orderStatusFlow[s]
	return next, ok
}

func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusAccepted, OrderStatusPreparing, OrderStatusOnTheWay, OrderStatusDelivered:
		return true
	default:
		return false
	}
}

type OrderItem struct {
	DishID   int64   `json:"dish_id"`
	DishName string  `json:"dish_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id"`
	RestaurantID int64       `json:"restaurant_id"`
	CourierID    *int64      `json:"courier_id,omitempty"`
	Status       OrderStatus `json:"status"`
	Items        []OrderItem `json:"items"`
	Subtotal     float64     `json:"subtotal"`
	DeliveryFee  float64     `json:"delivery_fee"`
	Total        float64     `json:"total"`
	Address      string      `json:"address"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
