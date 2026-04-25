package controllers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lab7-oop/clients/sqlite"
	"lab7-oop/controllers"
	apperrors "lab7-oop/shared/errors"
)

func TestCartAddAndTotal(t *testing.T) {
	db := setupTestDB(t)
	menuRepo := sqlite.NewMenuRepository(db)
	cartRepo := sqlite.NewCartRepository(db)
	ctrl := controllers.NewCartController(cartRepo, menuRepo)
	ctx := context.Background()
	userID := int64(1)

	cart, err := ctrl.AddItem(ctx, userID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(1), cart.RestaurantID)
	assert.InDelta(t, 8.99*2, cart.Subtotal, 0.01)
	assert.InDelta(t, controllers.DeliveryFee, cart.DeliveryFee, 0.01)
	assert.InDelta(t, cart.Subtotal+cart.DeliveryFee, cart.Total, 0.01)

	_, err = ctrl.AddItem(ctx, userID, 4, 1)
	assert.ErrorIs(t, err, apperrors.ErrRestaurantMismatch)
}

func TestCartEmptyOrderBlocked(t *testing.T) {
	db := setupTestDB(t)
	menuRepo := sqlite.NewMenuRepository(db)
	cartRepo := sqlite.NewCartRepository(db)
	cartCtrl := controllers.NewCartController(cartRepo, menuRepo)
	orderRepo := sqlite.NewOrderRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)
	orderCtrl := controllers.NewOrderController(orderRepo, cartCtrl, courierRepo)

	_, err := orderCtrl.CreateFromCart(context.Background(), 1, "addr")
	assert.ErrorIs(t, err, apperrors.ErrEmptyCart)
}
