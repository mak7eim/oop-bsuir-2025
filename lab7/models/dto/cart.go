package dto

type AddToCartRequest struct {
	DishID   int64 `json:"dish_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	Address string `json:"address" binding:"required"`
}
