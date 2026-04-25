package handlers

import (
	"github.com/gin-gonic/gin"

	"lab7-oop/controllers"
	"lab7-oop/shared/response"
)

type CourierHandler struct {
	ctrl *controllers.CourierController
}

func NewCourierHandler(ctrl *controllers.CourierController) *CourierHandler {
	return &CourierHandler{ctrl: ctrl}
}

// ListCouriers godoc
// @Summary      List couriers
// @Tags         couriers
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Envelope
// @Router       /couriers [get]
func (h *CourierHandler) ListCouriers(c *gin.Context) {
	list, err := h.ctrl.List(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, list)
}
