package controllers

import (
	"context"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

const DeliveryFee = 3.99

type CartController struct {
	cart CartRepository
	menu MenuRepository
}

func NewCartController(cart CartRepository, menu MenuRepository) *CartController {
	return &CartController{cart: cart, menu: menu}
}

func (c *CartController) GetCart(ctx context.Context, userID int64) (*models.Cart, error) {
	return c.buildCart(ctx, userID)
}

func (c *CartController) AddItem(ctx context.Context, userID, dishID int64, quantity int) (*models.Cart, error) {
	dish, err := c.menu.GetDish(ctx, dishID)
	if err != nil {
		return nil, err
	}
	if !dish.Available {
		return nil, apperrors.ErrBadRequest
	}

	items, err := c.cart.GetItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(items) > 0 {
		cart, err := c.buildCart(ctx, userID)
		if err != nil {
			return nil, err
		}
		if cart.RestaurantID != 0 && cart.RestaurantID != dish.RestaurantID {
			return nil, apperrors.ErrRestaurantMismatch
		}
	}

	currentQty := items[dishID]
	if err := c.cart.SetItem(ctx, userID, dishID, currentQty+quantity); err != nil {
		return nil, err
	}
	return c.buildCart(ctx, userID)
}

func (c *CartController) RemoveItem(ctx context.Context, userID, dishID int64) (*models.Cart, error) {
	if err := c.cart.RemoveItem(ctx, userID, dishID); err != nil {
		return nil, err
	}
	return c.buildCart(ctx, userID)
}

func (c *CartController) Clear(ctx context.Context, userID int64) error {
	return c.cart.Clear(ctx, userID)
}

func (c *CartController) buildCart(ctx context.Context, userID int64) (*models.Cart, error) {
	items, err := c.cart.GetItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart := &models.Cart{
		UserID:      userID,
		Items:       []models.CartItem{},
		DeliveryFee: DeliveryFee,
	}

	for dishID, qty := range items {
		dish, err := c.menu.GetDish(ctx, dishID)
		if err != nil {
			return nil, err
		}
		if cart.RestaurantID == 0 {
			cart.RestaurantID = dish.RestaurantID
		}
		lineTotal := dish.Price * float64(qty)
		cart.Subtotal += lineTotal
		cart.Items = append(cart.Items, models.CartItem{
			DishID:   dish.ID,
			DishName: dish.Name,
			Quantity: qty,
			Price:    dish.Price,
		})
	}

	if len(cart.Items) == 0 {
		cart.DeliveryFee = 0
	}
	cart.Total = cart.Subtotal + cart.DeliveryFee
	return cart, nil
}
