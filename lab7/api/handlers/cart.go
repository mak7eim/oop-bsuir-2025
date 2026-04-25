package handlers

import (
	"github.com/gin-gonic/gin"

	"lab7-oop/api/middleware"
	"lab7-oop/controllers"
	"lab7-oop/models/dto"
	"lab7-oop/shared/response"
)

type CartHandler struct {
	ctrl *controllers.CartController
}

func NewCartHandler(ctrl *controllers.CartController) *CartHandler {
	return &CartHandler{ctrl: ctrl}
}

// GetCart godoc
// @Summary      Get current user's cart
// @Tags         cart
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	cart, err := h.ctrl.GetCart(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, cart)
}

// AddToCart godoc
// @Summary      Add dish to cart
// @Tags         cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.AddToCartRequest  true  "Item"
// @Success      200   {object}  response.Envelope
// @Router       /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	cart, err := h.ctrl.AddItem(c.Request.Context(), middleware.UserID(c), req.DishID, req.Quantity)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, cart)
}

// RemoveFromCart godoc
// @Summary      Remove dish from cart
// @Tags         cart
// @Security     BearerAuth
// @Param        dishId  path  int  true  "Dish ID"
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /cart/items/{dishId} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	dishID, err := parseID(c.Param("dishId"))
	if err != nil {
		response.Error(c, err)
		return
	}
	cart, err := h.ctrl.RemoveItem(c.Request.Context(), middleware.UserID(c), dishID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, cart)
}
