package controllers

import (
	"context"
	"time"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

type OrderController struct {
	orders  OrderRepository
	cart    *CartController
	couriers CourierRepository
}

func NewOrderController(orders OrderRepository, cart *CartController, couriers CourierRepository) *OrderController {
	return &OrderController{orders: orders, cart: cart, couriers: couriers}
}

func (c *OrderController) CreateFromCart(ctx context.Context, userID int64, address string) (*models.Order, error) {
	cart, err := c.cart.buildCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, apperrors.ErrEmptyCart
	}

	order := &models.Order{
		UserID:       userID,
		RestaurantID: cart.RestaurantID,
		Status:       models.OrderStatusAccepted,
		Subtotal:     cart.Subtotal,
		DeliveryFee:  cart.DeliveryFee,
		Total:        cart.Total,
		Address:      address,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	for _, item := range cart.Items {
		order.Items = append(order.Items, models.OrderItem{
			DishID:   item.DishID,
			DishName: item.DishName,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	if err := c.orders.Create(ctx, order); err != nil {
		return nil, err
	}
	_ = c.cart.Clear(ctx, userID)
	return order, nil
}

func (c *OrderController) GetByID(ctx context.Context, userID int64, role models.Role, orderID int64) (*models.Order, error) {
	order, err := c.orders.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if role == models.RoleCustomer && order.UserID != userID {
		return nil, apperrors.ErrForbidden
	}
	return order, nil
}

func (c *OrderController) List(ctx context.Context, userID int64, role models.Role) ([]models.Order, error) {
	if role == models.RoleAdmin {
		return c.orders.ListAll(ctx)
	}
	return c.orders.ListByUser(ctx, userID)
}

func (c *OrderController) AdvanceStatus(ctx context.Context, orderID int64) (*models.Order, error) {
	order, err := c.orders.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	next, ok := order.Status.Next()
	if !ok {
		return nil, apperrors.ErrInvalidStatus
	}
	order.Status = next
	if err := c.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	if order.Status == models.OrderStatusDelivered && order.CourierID != nil {
		_ = c.couriers.UpdateStatus(ctx, *order.CourierID, models.CourierStatusFree)
	}
	return order, nil
}

func (c *OrderController) SetStatus(ctx context.Context, orderID int64, status models.OrderStatus) (*models.Order, error) {
	if !status.Valid() {
		return nil, apperrors.ErrBadRequest
	}
	order, err := c.orders.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.Status = status
	if err := c.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	if status == models.OrderStatusDelivered && order.CourierID != nil {
		_ = c.couriers.UpdateStatus(ctx, *order.CourierID, models.CourierStatusFree)
	}
	return order, nil
}

func (c *OrderController) AssignCourier(ctx context.Context, orderID int64, courierID *int64, auto bool) (*models.Order, error) {
	order, err := c.orders.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.CourierID != nil {
		return nil, apperrors.ErrConflict
	}

	var courier *models.Courier
	if auto {
		courier, err = c.couriers.FindFree(ctx)
	} else {
		if courierID == nil {
			return nil, apperrors.ErrBadRequest
		}
		courier, err = c.couriers.GetByID(ctx, *courierID)
		if err != nil {
			return nil, err
		}
		if courier.Status != models.CourierStatusFree {
			return nil, apperrors.ErrCourierBusy
		}
	}
	if err != nil {
		return nil, err
	}

	order.CourierID = &courier.ID
	if err := c.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	if err := c.couriers.UpdateStatus(ctx, courier.ID, models.CourierStatusBusy); err != nil {
		return nil, err
	}
	return order, nil
}
