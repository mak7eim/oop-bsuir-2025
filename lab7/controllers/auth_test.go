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
	apperrors "lab7-oop/shared/errors"
)

func setupTestDB(t *testing.T) *sqlite.DB {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, db.SeedDemoData(context.Background()))
	return db
}

func TestAuthRegisterAndLogin(t *testing.T) {
	db := setupTestDB(t)
	userRepo := sqlite.NewUserRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)
	hasher := auth.BcryptHasher{}
	tokenSvc := auth.NewJWTService("test-secret", 0)
	ctrl := controllers.NewAuthController(userRepo, tokenSvc, hasher, courierRepo)

	ctx := context.Background()
	reg, err := ctrl.Register(ctx, dto.RegisterRequest{
		Email:    "user@test.com",
		Password: "secret12",
		Role:     models.RoleCustomer,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, reg.Token)
	assert.Equal(t, "user@test.com", reg.User.Email)

	_, err = ctrl.Register(ctx, dto.RegisterRequest{
		Email:    "user@test.com",
		Password: "secret12",
		Role:     models.RoleCustomer,
	})
	assert.ErrorIs(t, err, apperrors.ErrConflict)

	login, err := ctrl.Login(ctx, dto.LoginRequest{Email: "user@test.com", Password: "secret12"})
	require.NoError(t, err)
	assert.NotEmpty(t, login.Token)

	_, err = ctrl.Login(ctx, dto.LoginRequest{Email: "user@test.com", Password: "wrong"})
	assert.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
}
