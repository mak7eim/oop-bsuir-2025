package models

type CourierStatus string

const (
	CourierStatusFree CourierStatus = "free"
	CourierStatusBusy CourierStatus = "busy"
)

type Courier struct {
	ID     int64         `json:"id"`
	UserID int64         `json:"user_id"`
	Name   string        `json:"name"`
	Phone  string        `json:"phone"`
	Status CourierStatus `json:"status"`
}
