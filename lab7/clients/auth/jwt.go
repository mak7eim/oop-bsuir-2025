package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

type claims struct {
	UserID int64       `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), ttl: ttl}
}

func (s *JWTService) Generate(userID int64, role models.Role) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})
	return token.SignedString(s.secret)
}

func (s *JWTService) Parse(tokenStr string) (int64, models.Role, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, "", apperrors.ErrUnauthorized
	}
	c, ok := token.Claims.(*claims)
	if !ok {
		return 0, "", apperrors.ErrUnauthorized
	}
	return c.UserID, c.Role, nil
}
