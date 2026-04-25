package dto

import "lab7-oop/models"

type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

type AssignCourierRequest struct {
	CourierID *int64 `json:"courier_id"`
	Auto      bool   `json:"auto"`
}
