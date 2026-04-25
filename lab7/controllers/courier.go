package controllers

import (
	"context"

	"lab7-oop/models"
)

type CourierController struct {
	couriers CourierRepository
	users    UserRepository
}

func NewCourierController(couriers CourierRepository, users UserRepository) *CourierController {
	return &CourierController{couriers: couriers, users: users}
}

func (c *CourierController) List(ctx context.Context) ([]models.Courier, error) {
	return c.couriers.List(ctx)
}

func (c *CourierController) EnsureProfile(ctx context.Context, userID int64, name, phone string) (*models.Courier, error) {
	list, err := c.couriers.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].UserID == userID {
			return &list[i], nil
		}
	}
	courier := &models.Courier{
		UserID: userID,
		Name:   name,
		Phone:  phone,
		Status: models.CourierStatusFree,
	}
	if err := c.couriers.Create(ctx, courier); err != nil {
		return nil, err
	}
	return courier, nil
}
