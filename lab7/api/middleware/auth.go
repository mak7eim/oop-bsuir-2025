package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"lab7-oop/controllers"
	"lab7-oop/models"
	apperrors "lab7-oop/shared/errors"
	"lab7-oop/shared/response"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

type AuthMiddleware struct {
	tokens controllers.TokenService
}

func NewAuthMiddleware(tokens controllers.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokens: tokens}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c.GetHeader("Authorization"))
		if token == "" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		userID, role, err := m.tokens.Parse(token)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, userID)
		c.Set(ContextRoleKey, role)
		c.Next()
	}
}

func (m *AuthMiddleware) RequireRoles(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(ContextRoleKey)
		if !ok {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		userRole := role.(models.Role)
		for _, allowed := range roles {
			if userRole == allowed {
				c.Next()
				return
			}
		}
		response.Error(c, apperrors.ErrForbidden)
		c.Abort()
	}
}

func UserID(c *gin.Context) int64 {
	return c.GetInt64(ContextUserIDKey)
}

func Role(c *gin.Context) models.Role {
	v, _ := c.Get(ContextRoleKey)
	return v.(models.Role)
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}
