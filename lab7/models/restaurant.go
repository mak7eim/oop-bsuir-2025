package models

type Restaurant struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
}

type Dish struct {
	ID           int64   `json:"id"`
	RestaurantID int64   `json:"restaurant_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
}
