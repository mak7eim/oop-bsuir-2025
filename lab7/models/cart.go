package models

type CartItem struct {
	DishID   int64   `json:"dish_id"`
	DishName string  `json:"dish_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Cart struct {
	UserID       int64      `json:"user_id"`
	RestaurantID int64      `json:"restaurant_id"`
	Items        []CartItem `json:"items"`
	Subtotal     float64    `json:"subtotal"`
	DeliveryFee  float64    `json:"delivery_fee"`
	Total        float64    `json:"total"`
}
