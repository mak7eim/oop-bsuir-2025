package controllers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lab7-oop/clients/auth"
	"lab7-oop/clients/sqlite"
	"lab7-oop/controllers"
	"lab7-oop/models"
	"lab7-oop/models/dto"
)

func TestOrderFlowAndCourier(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	userRepo := sqlite.NewUserRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)
	menuRepo := sqlite.NewMenuRepository(db)
	cartRepo := sqlite.NewCartRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)

	hasher := auth.BcryptHasher{}
	hash, _ := hasher.Hash("pass")
	require.NoError(t, userRepo.Create(ctx, &models.User{Email: "c@test.com", PasswordHash: hash, Role: models.RoleCustomer}))
	require.NoError(t, userRepo.Create(ctx, &models.User{Email: "cr@test.com", PasswordHash: hash, Role: models.RoleCourier}))

	couriers, _ := courierRepo.List(ctx)
	var courierID int64
	if len(couriers) == 0 {
		require.NoError(t, courierRepo.Create(ctx, &models.Courier{UserID: 2, Name: "Courier", Status: models.CourierStatusFree}))
		couriers, _ = courierRepo.List(ctx)
	}
	courierID = couriers[0].ID

	cartCtrl := controllers.NewCartController(cartRepo, menuRepo)
	orderCtrl := controllers.NewOrderController(orderRepo, cartCtrl, courierRepo)

	_, err := cartCtrl.AddItem(ctx, 1, 1, 1)
	require.NoError(t, err)

	order, err := orderCtrl.CreateFromCart(ctx, 1, "Home 1")
	require.NoError(t, err)
	assert.Equal(t, models.OrderStatusAccepted, order.Status)

	order, err = orderCtrl.AdvanceStatus(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.OrderStatusPreparing, order.Status)

	order, err = orderCtrl.AssignCourier(ctx, order.ID, &courierID, false)
	require.NoError(t, err)
	assert.Equal(t, courierID, *order.CourierID)

	courier, _ := courierRepo.GetByID(ctx, courierID)
	assert.Equal(t, models.CourierStatusBusy, courier.Status)

	order, err = orderCtrl.AdvanceStatus(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.OrderStatusOnTheWay, order.Status)

	order, err = orderCtrl.AdvanceStatus(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.OrderStatusDelivered, order.Status)

	courier, _ = courierRepo.GetByID(ctx, courierID)
	assert.Equal(t, models.CourierStatusFree, courier.Status)
}

func TestAutoAssignCourier(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	menuRepo := sqlite.NewMenuRepository(db)
	cartRepo := sqlite.NewCartRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)
	require.NoError(t, courierRepo.Create(ctx, &models.Courier{UserID: 99, Name: "Free", Status: models.CourierStatusFree}))

	cartCtrl := controllers.NewCartController(cartRepo, menuRepo)
	orderCtrl := controllers.NewOrderController(orderRepo, cartCtrl, courierRepo)

	_, _ = cartCtrl.AddItem(ctx, 1, 2, 1)
	order, err := orderCtrl.CreateFromCart(ctx, 1, "Addr")
	require.NoError(t, err)

	order, err = orderCtrl.AssignCourier(ctx, order.ID, nil, true)
	require.NoError(t, err)
	assert.NotNil(t, order.CourierID)
}

func TestAuthControllerCourierRegister(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()
	userRepo := sqlite.NewUserRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)
	hasher := auth.BcryptHasher{}
	tokenSvc := auth.NewJWTService("secret", 0)
	ctrl := controllers.NewAuthController(userRepo, tokenSvc, hasher, courierRepo)

	_, err := ctrl.Register(ctx, dto.RegisterRequest{
		Email: "courier@test.com", Password: "secret12", Role: models.RoleCourier,
	})
	require.NoError(t, err)

	list, err := courierRepo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, models.CourierStatusFree, list[0].Status)
}
