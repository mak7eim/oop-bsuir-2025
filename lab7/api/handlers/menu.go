package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"lab7-oop/controllers"
	"lab7-oop/shared/response"
)

type MenuHandler struct {
	ctrl *controllers.MenuController
}

func NewMenuHandler(ctrl *controllers.MenuController) *MenuHandler {
	return &MenuHandler{ctrl: ctrl}
}

// ListRestaurants godoc
// @Summary      List restaurants
// @Tags         menu
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /restaurants [get]
func (h *MenuHandler) ListRestaurants(c *gin.Context) {
	list, err := h.ctrl.ListRestaurants(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}

// GetRestaurant godoc
// @Summary      Get restaurant by ID
// @Tags         menu
// @Produce      json
// @Param        id  path  int  true  "Restaurant ID"
// @Success      200  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Router       /restaurants/{id} [get]
func (h *MenuHandler) GetRestaurant(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	rest, err := h.ctrl.GetRestaurant(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, rest)
}

// ListDishes godoc
// @Summary      List dishes of a restaurant
// @Tags         menu
// @Produce      json
// @Param        id  path  int  true  "Restaurant ID"
// @Success      200  {object}  response.Envelope
// @Router       /restaurants/{id}/dishes [get]
func (h *MenuHandler) ListDishes(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	dishes, err := h.ctrl.ListDishes(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, dishes)
}

func parseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
