package handlers

import (
	"github.com/gin-gonic/gin"

	"lab7-oop/api/middleware"
	"lab7-oop/controllers"
	"lab7-oop/models"
	"lab7-oop/models/dto"
	"lab7-oop/shared/response"
)

type OrderHandler struct {
	ctrl *controllers.OrderController
}

func NewOrderHandler(ctrl *controllers.OrderController) *OrderHandler {
	return &OrderHandler{ctrl: ctrl}
}

// CreateOrder godoc
// @Summary      Create order from cart
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.CreateOrderRequest  true  "Order"
// @Success      201   {object}  response.Envelope
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	order, err := h.ctrl.CreateFromCart(c.Request.Context(), middleware.UserID(c), req.Address)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, order)
}

// ListOrders godoc
// @Summary      List orders
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.ctrl.List(c.Request.Context(), middleware.UserID(c), middleware.Role(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, orders)
}

// GetOrder godoc
// @Summary      Get order by ID
// @Tags         orders
// @Security     BearerAuth
// @Param        id  path  int  true  "Order ID"
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	order, err := h.ctrl.GetByID(c.Request.Context(), middleware.UserID(c), middleware.Role(c), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, order)
}

// AdvanceOrderStatus godoc
// @Summary      Advance order to next status
// @Description  accepted -> preparing -> on_the_way -> delivered
// @Tags         orders
// @Security     BearerAuth
// @Param        id  path  int  true  "Order ID"
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /orders/{id}/advance [post]
func (h *OrderHandler) AdvanceOrderStatus(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	order, err := h.ctrl.AdvanceStatus(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, order)
}

// SetOrderStatus godoc
// @Summary      Set order status (admin)
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Order ID"
// @Param        body  body  dto.UpdateOrderStatusRequest  true  "Status"
// @Success      200   {object}  response.Envelope
// @Router       /orders/{id}/status [put]
func (h *OrderHandler) SetOrderStatus(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	order, err := h.ctrl.SetStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, order)
}

// AssignCourier godoc
// @Summary      Assign courier to order (admin)
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Order ID"
// @Param        body  body  dto.AssignCourierRequest  true  "Assignment"
// @Success      200   {object}  response.Envelope
// @Router       /orders/{id}/courier [post]
func (h *OrderHandler) AssignCourier(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.AssignCourierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	order, err := h.ctrl.AssignCourier(c.Request.Context(), id, req.CourierID, req.Auto)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, order)
}

func RequireCustomerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := middleware.Role(c)
		if role != models.RoleCustomer && role != models.RoleAdmin {
			c.Next()
			return
		}
		c.Next()
	}
}
