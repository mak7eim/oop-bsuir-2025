package controllers

import (
	"context"

	"lab7-oop/models"
	"lab7-oop/models/dto"
	apperrors "lab7-oop/shared/errors"
)

type AuthController struct {
	users    UserRepository
	tokens   TokenService
	password PasswordHasher
	couriers CourierRepository
}

func NewAuthController(users UserRepository, tokens TokenService, password PasswordHasher, couriers CourierRepository) *AuthController {
	return &AuthController{users: users, tokens: tokens, password: password, couriers: couriers}
}

func (c *AuthController) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if req.Role != models.RoleCustomer && req.Role != models.RoleAdmin && req.Role != models.RoleCourier {
		return nil, apperrors.ErrBadRequest
	}

	hash, err := c.password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: req.Email, PasswordHash: hash, Role: req.Role}
	if err := c.users.Create(ctx, user); err != nil {
		return nil, err
	}

	if req.Role == models.RoleCourier && c.couriers != nil {
		_ = c.couriers.Create(ctx, &models.Courier{
			UserID: user.ID,
			Name:   user.Email,
			Status: models.CourierStatusFree,
		})
	}

	token, err := c.tokens.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (c *AuthController) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := c.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}
	if err := c.password.Compare(user.PasswordHash, req.Password); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := c.tokens.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (c *AuthController) Me(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := c.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := toUserResponse(user)
	return &resp, nil
}

func toUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{ID: u.ID, Email: u.Email, Role: u.Role}
}
