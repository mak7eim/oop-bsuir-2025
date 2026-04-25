package errors

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrConflict          = errors.New("resource already exists")
	ErrBadRequest        = errors.New("bad request")
	ErrInvalidStatus     = errors.New("invalid order status transition")
	ErrCourierBusy       = errors.New("courier is not available")
	ErrEmptyCart         = errors.New("cart is empty")
	ErrRestaurantMismatch = errors.New("all items must be from the same restaurant")
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAppError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func (e *AppError) Error() string {
	return e.Message
}

func From(err error) *AppError {
	switch {
	case errors.Is(err, ErrNotFound):
		return NewAppError("not_found", err.Error())
	case errors.Is(err, ErrUnauthorized):
		return NewAppError("unauthorized", err.Error())
	case errors.Is(err, ErrForbidden):
		return NewAppError("forbidden", err.Error())
	case errors.Is(err, ErrInvalidCredentials):
		return NewAppError("invalid_credentials", err.Error())
	case errors.Is(err, ErrConflict):
		return NewAppError("conflict", err.Error())
	case errors.Is(err, ErrInvalidStatus):
		return NewAppError("invalid_status", err.Error())
	case errors.Is(err, ErrCourierBusy):
		return NewAppError("courier_busy", err.Error())
	case errors.Is(err, ErrEmptyCart):
		return NewAppError("empty_cart", err.Error())
	case errors.Is(err, ErrRestaurantMismatch):
		return NewAppError("restaurant_mismatch", err.Error())
	default:
		return NewAppError("bad_request", err.Error())
	}
}
