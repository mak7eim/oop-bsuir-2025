package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"lab7-oop/api/middleware"
	"lab7-oop/controllers"
	"lab7-oop/models/dto"
	"lab7-oop/shared/response"
)

type AuthHandler struct {
	ctrl *controllers.AuthController
}

func NewAuthHandler(ctrl *controllers.AuthController) *AuthHandler {
	return &AuthHandler{ctrl: ctrl}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.RegisterRequest  true  "Registration data"
// @Success      201   {object}  response.Envelope{data=dto.AuthResponse}
// @Failure      400   {object}  response.Envelope
// @Failure      409   {object}  response.Envelope
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.ctrl.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.LoginRequest  true  "Credentials"
// @Success      200   {object}  response.Envelope{data=dto.AuthResponse}
// @Failure      401   {object}  response.Envelope
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.ctrl.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// Me godoc
// @Summary      Get current user
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      401  {object}  response.Envelope
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	res, err := h.ctrl.Me(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// Health godoc
// @Summary      Health check
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
