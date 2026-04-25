package controllers

import (
	"context"

	"lab7-oop/models"
)

type MenuController struct {
	menu MenuRepository
}

func NewMenuController(menu MenuRepository) *MenuController {
	return &MenuController{menu: menu}
}

func (c *MenuController) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	return c.menu.ListRestaurants(ctx)
}

func (c *MenuController) GetRestaurant(ctx context.Context, id int64) (*models.Restaurant, error) {
	return c.menu.GetRestaurant(ctx, id)
}

func (c *MenuController) ListDishes(ctx context.Context, restaurantID int64) ([]models.Dish, error) {
	if _, err := c.menu.GetRestaurant(ctx, restaurantID); err != nil {
		return nil, err
	}
	return c.menu.ListDishesByRestaurant(ctx, restaurantID)
}
