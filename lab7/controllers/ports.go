package controllers

import (
	"context"

	"lab7-oop/models"
)

// Repository interfaces (Dependency Inversion): controllers depend on abstractions,
// clients provide concrete implementations.

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type MenuRepository interface {
	ListRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetRestaurant(ctx context.Context, id int64) (*models.Restaurant, error)
	ListDishesByRestaurant(ctx context.Context, restaurantID int64) ([]models.Dish, error)
	GetDish(ctx context.Context, id int64) (*models.Dish, error)
}

type CartRepository interface {
	GetItems(ctx context.Context, userID int64) (map[int64]int, error)
	SetItem(ctx context.Context, userID, dishID int64, quantity int) error
	RemoveItem(ctx context.Context, userID, dishID int64) error
	Clear(ctx context.Context, userID int64) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	ListByUser(ctx context.Context, userID int64) ([]models.Order, error)
	ListAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order *models.Order) error
}

type CourierRepository interface {
	List(ctx context.Context) ([]models.Courier, error)
	GetByID(ctx context.Context, id int64) (*models.Courier, error)
	FindFree(ctx context.Context) (*models.Courier, error)
	UpdateStatus(ctx context.Context, id int64, status models.CourierStatus) error
	Create(ctx context.Context, courier *models.Courier) error
}

type TokenService interface {
	Generate(userID int64, role models.Role) (string, error)
	Parse(token string) (userID int64, role models.Role, err error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
